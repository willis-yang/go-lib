package utils

import (
	"context"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestGetAccessIp(t *testing.T) {
	// 创建一个带有元数据的上下文
	md := metadata.New(map[string]string{"X-Forwarded-For": "192.168.1.1"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	ipAddr, err := GetMetaClientIp(ctx)
	if err != nil {
		t.Errorf("GetAccessIp failed with error %v", err)
		return
	}
	t.Logf("Client IP: %s", ipAddr)

}
