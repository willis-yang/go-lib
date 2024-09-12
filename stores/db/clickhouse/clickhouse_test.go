package clickhouse

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestNewClickHouse(t *testing.T) {
	conf := ClickHouseConfig{
		Adders:          []string{"192.168.12.21:18123"},
		Databases:       "common",
		Username:        "default",
		Password:        "",
		DialTimeout:     10,
		Tracing:         false,
		MaxIdle:         1,
		MaxOpen:         1,
		ConnMaxIdleTime: 5,
		Debug:           true,
		TLS:             false,
		ConnMaxLifeTime: 10,
	}
	clickhouseConn := NewClickHouse(conf, logx.LogConf{
		ServiceName:         "",
		Mode:                "",
		Encoding:            "",
		TimeFormat:          "",
		Path:                "",
		Level:               "",
		MaxContentLength:    0,
		Compress:            false,
		Stat:                false,
		KeepDays:            0,
		StackCooldownMillis: 0,
		MaxBackups:          0,
		MaxSize:             0,
		Rotation:            "",
	})
	err := clickhouseConn.Exec("select count(*) from common.country_list").Error
	if err != nil {
		panic(err)
	}
}
