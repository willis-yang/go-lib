package dingdingTalk

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestNewWebhook(t *testing.T) {
	dingTalk := NewWebhook(
		"xxxxx",
		"xxxx")
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     "127.0.0.1:6379",
			DB:       0,
			Password: "",
		})
	dingTalk.SetRedis(redisClient)
	dingTalk.SetFrequencyLimit(Limit{
		InUse: true,
		Limit: 1,
	})
	err := dingTalk.SendTextMessage(context.Background(), TextMessage{
		Text: struct {
			Content string `json:"content"`
		}{
			Content: "test msg",
		},
	})
	if err != nil {
		t.Error(err)
	}
}
