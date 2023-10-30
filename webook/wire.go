// //go:build wireinject
package main

//
//import (
//	"basic-go/webook/internal/repository"
//	"basic-go/webook/internal/repository/cache"
//	"basic-go/webook/internal/repository/dao"
//	"basic-go/webook/internal/service"
//	"basic-go/webook/internal/web"
//	"basic-go/webook/ioc"
//	"github.com/gin-gonic/gin"
//	"github.com/google/wire"
//)
//
//func InitWebServer() *gin.Engine {
//	wire.Build(ioc.InitRedis, ioc.InitDB,
//		dao.NewGORMUserDao,
//		cache.NewRedisUserCache, cache.NewRedisCodeCache,
//		repository.NewCacheUserRepository, repository.NewCacheCodeRepository,
//
//		ioc.InitSMSService,
//		service.NewuserService, service.NewCodeRedisService,
//
//		web.NewUserHandler,
//
//		ioc.InitGinMiddlewares,
//
//		ioc.InitWebServer,
//	)
//	return gin.Default()
//}
