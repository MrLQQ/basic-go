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
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDAO, userCache)
	userService := service.NewuserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	articleRedisCache := cache.NewArticleRedisCache(cmdable)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	codeRepository := repository.NewCacheCodeRepository(codeCache)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, articleRedisCache)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache)
	smsService := ioc.InitSMSService(logger)
	wechatService := InitWechatService(logger)
	codeService := service.NewCodeCacheService(codeRepository, smsService)
	articleService := service.NewArticleService(articleRepository)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	userHandler := web.NewUserHandler(userService, handler, codeService, logger)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, logger)
	wechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService, logger)
	engine := ioc.InitWebServer(v, userHandler, wechatHandler, articleHandler)
	return engine
}

func InitArticleHandler(Dao dao.ArticleDAO) *web.ArticleHandler {
	cmdable := ioc.InitRedis()
	loggerV1 := InitLogger()
	db := InitDB()
	//articleDAO := dao.NewArticleGORMDAO(db)
	articleRedisCache := cache.NewArticleRedisCache(cmdable)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	articleDAO := Dao
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, articleRedisCache)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache)
	articleService := service.NewArticleService(articleRepository)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, loggerV1)
	return articleHandler
}

func InitInteractiveService() service.InteractiveService {
	db := InitDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	return interactiveService
}
