package main

import (
	"basic-go/webook/internal/events/article"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/ioc"
)

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	logger := ioc.InitLogger()
	handler := ijwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitGinMiddlewares(cmdable, handler, logger)
	db := ioc.InitDB(logger)
	userDAO := dao.NewGORMUserDao(db)
	articleDAO := dao.NewArticleGORMDAO(db)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache)
	userService := service.NewuserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCacheCodeRepository(codeCache)
	articleRedisCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, articleRedisCache)
	smsService := ioc.InitSMSService(logger)
	wechatService := ioc.InitWechatService(logger)
	codeService := service.NewCodeCacheService(codeRepository, smsService)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	v2 := ioc.InitConsumers()
	interactiveService := service.NewInteractiveService(interactiveRepository)
	userHandler := web.NewUserHandler(userService, handler, codeService, logger)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, logger)
	wechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService, logger)
	engine := ioc.InitWebServer(v, userHandler, wechatHandler, articleHandler)
	app := &App{
		server:    engine,
		consumers: v2,
	}
	return app
}
