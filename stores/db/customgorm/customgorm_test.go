package customgorm

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
	"time"
)

func TestNewGorm(t *testing.T) {
	gorm := NewGorm(GormConfig{
		DSN:            "root:xxxxx@tcp(xxxxx:3306)/release_manage?charset=utf8mb4&parseTime=true&loc=Local&timeout=3600s",
		TablePrefix:    "",
		DatabaseType:   DatabaseTypeMysql,
		MaxIdLe:        0,
		MaxConnect:     0,
		ConnectMaxLife: 0,
	}, GormLogConfig{
		Path:                      "",
		Level:                     "error",
		SlowThreshold:             time.Second,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
		Colorful:                  true,
	})
	var list []string
	err := gorm.Select("id").Where("id = ?", 10).Table("users").Find(&list).Error
	if err != nil {
		panic(err)
	}
	logx.Info(list)
}
