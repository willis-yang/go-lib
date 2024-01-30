package clickhouse

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestNewClickHouse(t *testing.T) {
	conf := ClickHouseConfig{
		Adders:          []string{"192.168.33.10:18123"},
		Databases:       "cloud_phone_statistic",
		Username:        "default",
		Password:        "11@11",
		DialTimeout:     0,
		Tracing:         false,
		MaxIdle:         0,
		MaxOpen:         0,
		ConnMaxIdleTime: 0,
		Debug:           true,
		TLS:             false,
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
	err := clickhouseConn.Exec("select count(*) from cloud_phone_statistic.cps_event_tracking").Error
	if err != nil {
		panic(err)
	}
}
