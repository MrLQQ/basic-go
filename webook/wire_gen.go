package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
)

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.InitGinMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewGORMUserDao(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewuserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCacheCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	wechatService := ioc.InitWechatService()
	codeService := service.NewCodeCacheService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService)
	wechatHandler := web.NewOAuth2WechatHandler(wechatService)
	engine := ioc.InitWebServer(v, userHandler, wechatHandler)
	return engine
}
