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
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "authorization"},
		//AllowMethods:     []string{"POST"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "myself.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	hdl.RegisterRoutes(server)

	server.Run(":8080")
}
