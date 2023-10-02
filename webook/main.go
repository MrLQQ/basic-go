package main

import (
	"basic-go/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	hdl := web.NewUserHandler()

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	hdl.RegisterRoutes(server)

	server.Run(":8080")
}
