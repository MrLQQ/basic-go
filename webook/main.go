package main

import (
	"basic-go/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func main() {
	hdl := web.NewUserHandler()

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//允许的请求源
		AllowOrigins: []string{"http://localhost:3000"},
		//是否允许带上用户认证信息（比如cookie）
		AllowCredentials: true,
		//业务请求中可以带上的头
		AllowHeaders: []string{"Content-Type", "authorization"},
		//AllowMethods:     []string{"POST"},
		// 哪些来源是允许的
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "myself.com")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这是我的middleware……")
	})
	hdl.RegisterRoutes(server)

	server.Run(":8080")
}
