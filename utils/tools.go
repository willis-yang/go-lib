package utils

import (
	"context"
	"google.golang.org/grpc/peer"
	"gorm.io/gorm/logger"
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

func GetLogLevel(level string) logger.LogLevel {
	logLevel := logger.Info
	switch level {
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "silent":
		logLevel = logger.Silent
	default:
		logLevel = logger.Info
	}
	return logLevel
}
