package gormConnect

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// 自定义数据库连接池
// 支持读写分离模式
// 查询缓存
type GormConfig struct {
	//dns 连接参数：https://github.com/go-sql-driver/mysql#dsn-data-source-name
	DSN            string // dsh连接 [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	TablePrefix    string // 表前缀
	DatabaseType   string // 连接驱动类型
	MaxIdLe        int    //空闲最大连接数
	MaxConnect     int    //最大连接数
	ConnectMaxLife int    //可复用连接最长时间
}

const DatabaseTypeMysql = "mysql"
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
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
	} else {
		filePath = fmt.Sprintf("%v/sql-%v.log", logConfig.Path, time.Now().Format("2006-01-02"))
		file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		newLogger = logger.New(
			log.New(file, "\r\n", log.LstdFlags), // 文件
			//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                      // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		//连接异常直接抛异常
		panic(err)
	}

	sqlDB, err := db.DB()
	// 设置空闲连接池中的最大连接数
	sqlDB.SetMaxIdleConns(gormConfig.MaxIdLe)

	//设置与数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(gormConfig.MaxConnect)

	// 设置可重复使用连接的最长时间
	sqlDB.SetConnMaxLifetime(time.Duration(gormConfig.ConnectMaxLife))

	return db
}
