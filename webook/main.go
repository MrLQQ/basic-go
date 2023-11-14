package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"

	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/spf13/viper/remote"
)

func main() {
	initViperV1()
	initLogger()
	//initViperWatch()
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,启动成功")
	})
	server.Run(":8080")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	// 当前工作目录的config目录
	viper.AddConfigPath("config")
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

// initViperWatch 监听配置变更
func initViperWatch() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	// 这一步之后 cfalg才有值
	pflag.Parse()

	//viper.Set("db.dsn", "root:root@tcp(localhost:13316)/webook")
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println(viper.GetString("test.key"))
	})
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	// 这一步之后 cfalg才有值
	pflag.Parse()

	//viper.Set("db.dsn", "root:root@tcp(localhost:13316)/webook")
	viper.SetConfigFile(*cfile)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("远程配置中心发生变更")
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			log.Println("watch", viper.GetString("test.key"))
			time.Sleep(time.Second * 3)
		}
	}()
}
