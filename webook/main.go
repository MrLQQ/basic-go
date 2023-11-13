package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	initViperV1()
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,启动成功")
	})
	server.Run(":8080")
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

func initViperV1() {
	cfalg := pflag.String("config", "config/config.yaml", "配置文件路径")
	// 这一步之后 cfalg才有值
	pflag.Parse()

	//viper.Set("db.dsn", "root:root@tcp(localhost:13316)/webook")
	viper.SetConfigFile(*cfalg)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}
