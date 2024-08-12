package customgorm

import (
	"fmt"
	"github.com/willis-yang/go-lib/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type GormConfig struct {
	DSN            string // dsh [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	TablePrefix    string
	DatabaseType   string
	MaxIdLe        int
	MaxConnect     int
	ConnectMaxLife int
}

const DatabaseTypeMysql = "mysql"
const DatabaseTypeSqlite = "sqlite"
const DatabaseTypePostgreSQL = "PostgreSQL"
const DatabaseTypeTiDB = "TiDB"

// 初始化gorm 连接池
func NewGorm(gormConfig GormConfig, logConfig logx.LogConf) *gorm.DB {

	var (
		filePath  string
		newLogger logger.Interface
	)
	if logConfig.Path == "" { //无配置文件类型，直接输出日志
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  utils.GetLogLevel(logConfig.Level),
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			},
		)
	} else {
		filePath = fmt.Sprintf("%v/sql-%v.log", logConfig.Path, time.Now().Format("2006-01-02"))
		file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		newLogger = logger.New(
			log.New(file, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  utils.GetLogLevel(logConfig.Level),
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			},
		)
	}

	var gormDialector gorm.Dialector

	if gormConfig.DatabaseType == DatabaseTypeMysql {
		gormDialector = mysql.New(mysql.Config{
			DSN:                       gormConfig.DSN,
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		})
	} else if gormConfig.DatabaseType == DatabaseTypeSqlite {
		gormDialector = sqlite.Open(gormConfig.DSN)
	}
	db, err := gorm.Open(gormDialector, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(gormConfig.MaxIdLe)
	sqlDB.SetMaxOpenConns(gormConfig.MaxConnect)
	sqlDB.SetConnMaxLifetime(time.Duration(int64(gormConfig.ConnectMaxLife)) * time.Second)

	return db
}
