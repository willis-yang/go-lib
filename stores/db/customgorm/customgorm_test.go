package customgorm

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestNewGorm(t *testing.T) {
	gorm := NewGorm(GormConfig{
		DSN:            "root:xxxxx@tcp(xxxxx:3306)/release_manage?charset=utf8mb4&parseTime=true&loc=Local&timeout=3600s",
		TablePrefix:    "",
		DatabaseType:   DatabaseTypeMysql,
		MaxIdLe:        0,
		MaxConnect:     0,
		ConnectMaxLife: 0,
	}, logx.LogConf{
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
	var list []string
	err := gorm.Select("id").Where("id = ?", 10).Table("users").Find(&list).Error
	if err != nil {
		panic(err)
	}
	logx.Info(list)
}
