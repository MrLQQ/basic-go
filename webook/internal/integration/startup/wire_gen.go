package startup

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
	logger := InitLogger()
	handler := ijwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitGinMiddlewares(cmdable, handler, logger)
	db := InitDB()
	userDAO := dao.NewGORMUserDao(db)
	articleDAO := dao.NewArticleGORMDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewuserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCacheCodeRepository(codeCache)
	articleRepository := repository.NewCachedArticleRepository(articleDAO)
	smsService := ioc.InitSMSService(logger)
	wechatService := InitWechatService(logger)
	codeService := service.NewCodeCacheService(codeRepository, smsService)
	articleService := service.NewArticleService(articleRepository)
	userHandler := web.NewUserHandler(userService, handler, codeService, logger)
	articleHandler := web.NewArticleHandler(articleService, logger)
	wechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService, logger)
	engine := ioc.InitWebServer(v, userHandler, wechatHandler, articleHandler)
	return engine
}

func InitArticleHandler() *web.ArticleHandler {
	loggerV1 := InitLogger()
	db := InitDB()
	articleDAO := dao.NewArticleGORMDAO(db)
	articleRepository := repository.NewCachedArticleRepository(articleDAO)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	return articleHandler
}
