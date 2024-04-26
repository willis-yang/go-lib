package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"testing"
	"time"
)

func TestRetryFunc(t *testing.T) {
	opt := RetryOptions{
		MaxRetries: 3,
		Delay:      time.Second * 1,
	}
	fn := func() error {
		b := big.NewInt(3)
		serialNumber, err := rand.Int(rand.Reader, b)
		if err != nil {
			return err
		}
		if serialNumber.Int64() > 1 {
			return errors.New("exec failed")
		}
		return nil
	}
	err := RetryFunc(context.Background(), fn, opt)
	if err != nil {
		t.Errorf("RetryFunc error :%v", err)
	}
}
