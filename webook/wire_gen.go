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
	logger := ioc.InitLogger()
	handler := ijwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitGinMiddlewares(cmdable, handler, logger)
	db := ioc.InitDB(logger)
	userDAO := dao.NewGORMUserDao(db)
	articleDAO := dao.NewArticleGORMDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewuserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCacheCodeRepository(codeCache)
	articleRedisCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, articleRedisCache)
	smsService := ioc.InitSMSService(logger)
	wechatService := ioc.InitWechatService(logger)
	codeService := service.NewCodeCacheService(codeRepository, smsService)
	articleService := service.NewArticleService(articleRepository)
	userHandler := web.NewUserHandler(userService, handler, codeService, logger)
	articleHandler := web.NewArticleHandler(articleService, logger)
	wechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService, logger)
	engine := ioc.InitWebServer(v, userHandler, wechatHandler, articleHandler)
	return engine
}
