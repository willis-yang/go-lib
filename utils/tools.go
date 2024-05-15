package utils

import (
	"context"
	"google.golang.org/grpc/peer"
	"net"
)

func GetRemoteIp(ctx context.Context) (string, error) {
	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}
	}
	return addr, nil
}
