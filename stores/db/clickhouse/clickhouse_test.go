package clickhouse

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestNewClickHouse(t *testing.T) {
	conf := ClickHouseConfig{
		Dsn:             "clickhouse://default:@192.168.33.10:19000/cloud_phone_statistic?dial_timeout=200ms&max_execution_time=60", //clickhouse://default:@127.0.0.1:9000/default?dial_timeout=200ms&max_execution_time=60
		Tracing:         false,
		MaxIdle:         50,
		MaxOpen:         500,
		ConnMaxIdleTime: 3600,
	}

	log := logx.LogConf{
		Mode:  "Console",
		Level: "info",
	}
	db := NewClickHouse(conf, log)
	var count int
	err := db.Raw("Select count(*) as total from cps_event_tracking").Scan(&count).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(count)
	return
}
