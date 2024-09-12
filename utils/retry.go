package utils

import (
	"context"
	"time"
)

type RetryOptions struct {
	MaxRetries int           // Maximum number of retries
	Delay      time.Duration // Delay between retries
}

func RetryFunc(ctx context.Context, fn func() error, opts RetryOptions) error {
	var err error
	for i := 0; i < opts.MaxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(opts.Delay):
		}
	}

	return err
}
