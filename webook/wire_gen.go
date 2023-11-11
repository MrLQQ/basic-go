package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
)

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	handler := ijwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitGinMiddlewares(cmdable, handler)
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
	userHandler := web.NewUserHandler(userService, handler, codeService)
	wechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService)
	engine := ioc.InitWebServer(v, userHandler, wechatHandler)
	return engine
}
