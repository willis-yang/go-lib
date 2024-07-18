package utils

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

func GetMetaClientIp(ctx context.Context) (string, error) {
	var ip string
	if ctx == nil {
		return "", errors.New("ctx is empty")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		clientIPs := md.Get("X-Forwarded-For")
		if len(clientIPs) > 0 {
			ip = clientIPs[0]
		} else {
			return "", errors.New("client is empty")
		}
	}
	return ip, nil
}
