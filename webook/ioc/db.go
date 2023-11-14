package ioc

import (
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"DSN"`
	}
	// 提前给默认值
	//var cfg = Config{DSN: "root:root@tcp(localhost:13316)/webook"}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	//db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询日志
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{Key: "args", Value: i})
}
