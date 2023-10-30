package main

import (
	"net/http"
	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,启动成功")
	})
	server.Run(":8080")
}
