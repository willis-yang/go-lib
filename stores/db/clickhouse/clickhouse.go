package clickhouse

import (
	"crypto/tls"
	"fmt"
	clickhousego "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type ClickHouseConfig struct {
	Adders          []string
	Databases       string
	Username        string
	Password        string
	DialTimeout     int64
	Tracing         bool
	MaxIdle         int
	MaxOpen         int
	ConnMaxIdleTime int
	Debug           bool
	TLS             bool
}

func NewClickHouse(clickHouseConfig ClickHouseConfig, config logx.LogConf) *gorm.DB {
	var (
		filePath  string
		newLogger logger.Interface
	)
	if config.Path == "" {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			},
		)
	} else {
		filePath = fmt.Sprintf("%v/sql-%v.log", config.Path, time.Now().Format("2006-01-02"))
		file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		newLogger = logger.New(
			log.New(file, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			},
		)
	}
	clickhouseOptions := &clickhousego.Options{
		Addr: clickHouseConfig.Adders,
		Auth: clickhousego.Auth{
			Database: clickHouseConfig.Databases,
			Username: clickHouseConfig.Username,
			Password: clickHouseConfig.Password,
		},
		Settings: clickhousego.Settings{
			"max_execution_time": 60,
		},
		Protocol:    clickhousego.HTTP, //now is just supported http connection
		DialTimeout: 5 * time.Second,
		Compression: &clickhousego.Compression{
			Method: clickhousego.CompressionLZ4,
		},
		Debug: clickHouseConfig.Debug,
	}
	if clickHouseConfig.TLS {
		clickhouseOptions.TLS = &tls.Config{InsecureSkipVerify: clickHouseConfig.TLS}
	}
	clickhouseConn := clickhousego.OpenDB(clickhouseOptions)
	db, err := gorm.Open(clickhouse.New(clickhouse.Config{
		Conn:                         clickhouseConn,
		DisableDatetimePrecision:     true,
		DontSupportRenameColumn:      true,
		DontSupportEmptyDefaultValue: false,
		SkipInitializeWithVersion:    false,
		DefaultGranularity:           3,
		DefaultCompression:           "LZ4",
		DefaultIndexType:             "minmax",
		DefaultTableEngineOpts:       "ENGINE=MergeTree() ORDER BY tuple()",
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(clickHouseConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(clickHouseConfig.MaxOpen)
	sqlDB.SetConnMaxIdleTime(time.Duration(int64(clickHouseConfig.ConnMaxIdleTime)) * time.Second)

	return db
}
