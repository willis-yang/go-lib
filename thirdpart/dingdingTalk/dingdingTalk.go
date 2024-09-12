package dingdingTalk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"strconv"
	"time"
)

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Webhook struct {
	WebhookURL     string `json:"webhookURL"`
	Secret         string `json:"secret"`
	frequencyLimit Limit
	redis          *redis.Client
}

type Limit struct {
	Limit int64 `json:"limit"`
	InUse bool  `json:"inUse",default:false`
}

type TextMessage struct {
	At struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	MsgType string `json:"msgtype"`
}

type response struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

const (
	JsonContentType              = "application/json; charset=utf-8"
	DefaultMessageFrequencyLimit = 20
	Text                         = "text"
)

func NewWebhook(webhookURL, secret string) *Webhook {
	return &Webhook{
		WebhookURL: webhookURL,
		Secret:     secret,
	}
}

func (t *Webhook) SendTextMessage(ctx context.Context, msg TextMessage) error {
	if err := t.checkSendMessageLimit(ctx); err != nil {
		return err
	}
	msg.MsgType = Text
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(t.getURL(), JsonContentType, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("response error: %s", string(body)))
	}

	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New(fmt.Sprintf("response error: %s", string(body)))
	}
	return err
}

func (t *Webhook) SetRedis(redisClient *redis.Client) {
	t.redis = redisClient
}

func (t *Webhook) SetFrequencyLimit(frequencyLimit Limit) {
	t.frequencyLimit = frequencyLimit
}

// if send message failed, do not decrement the limit, wait for next minute
func (t *Webhook) checkSendMessageLimit(ctx context.Context) error {
	if t.frequencyLimit.InUse == false {
		return nil
	}
	if t.redis == nil {
		return errors.New(fmt.Sprintf("redis client is nil"))
	}
	key, err := t.getLimitKey()
	if err != nil {
		return err
	}
	res, err := t.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		_, err = t.redis.Set(ctx, key, 0, time.Minute).Result()
		if err != nil {
			return err
		}
		res = "0"
	}
	limit, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		return err
	}
	if limit >= t.frequencyLimit.Limit {
		return errors.New("exceeded sending limit")
	}
	_, err = t.redis.Incr(ctx, key).Result()
	return err
}

//func (t *Webhook) decrSendMessageLimit(ctx context.Context) error {
//	if t.frequencyLimit.InUse == false {
//		return nil
//	}
//	if t.redis == nil {
//		return nil
//	}
//	key, err := t.getLimitKey()
//	if err != nil {
//		return err
//	}
//	return t.redis.Decr(ctx, key).Err()
//}

func (t *Webhook) getLimitKey() (string, error) {
	md5Hash := md5.New()
	_, err := md5Hash.Write([]byte(t.WebhookURL))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("webhook:%ssendcount", fmt.Sprintf("%x", md5Hash.Sum(nil))), nil
}

func (t *Webhook) hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (t *Webhook) getURL() string {
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, t.Secret)
	sign := t.hmacSha256(stringToSign, t.Secret)
	return fmt.Sprintf("%s&timestamp=%d&sign=%s", t.WebhookURL, timestamp, sign)
}
