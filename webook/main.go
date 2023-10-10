package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {

	db := initDB()
	server := initWebServer()
	initUser(db, server)
	server.Run(":8080")
}

func initUser(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
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

	login := &middleware.LoginMiddlewareBuilder{}
	// 存储数据的，也就是userId存储在哪里
	// 直接存cookie
	store := cookie.NewStore([]byte("secret"))
	// 两个middleware,一个用来初始化session，一个用来登录校验
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}
