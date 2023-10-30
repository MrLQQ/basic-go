package main

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/localsms"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"github.com/coocood/freecache"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"strings"
	"time"

	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	// redis缓存
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	// 本地缓存
	//freeCache := freecache.NewCache(100 * 1024 * 1024)
	db := initDB()
	server := initWebServer()
	// 使用redis缓存处理验证码消息
	codeSvc := initCodeRedisSvc(redisClient)
	// 使用本地缓存处理验证码消息
	//codeSvc := initCodeSvc(freeCache)
	initUser(db, redisClient, codeSvc, server)
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,启动成功")
	})
	server.Run(":8080")
}

func initUser(db *gorm.DB, redisClient redis.Cmdable, codeSvc service.CodeService, server *gin.Engine) {
	ud := dao.NewGORMUserDao(db)
	userCache := cache.NewUserCache(redisClient)
	ur := repository.NewCacheUserRepository(ud, userCache)
	us := service.NewuserService(ur)
	hdl := web.NewUserHandler(us, codeSvc)
	hdl.RegisterRoutes(server)
}

/*
*
使用redis缓存实现
*/
func initCodeRedisSvc(redisClient redis.Cmdable) service.CodeService {
	crc := cache.NewRedisCodeCache(redisClient)
	crepo := repository.NewCacheCodeRepository(crc)
	return service.NewCodeRedisService(crepo, initMemorySms())
}

/*
*
使用本地缓存实现
*/
func initCodeSvc(freeCache *freecache.Cache) service.CodeService {
	cc := cache.NewMemoryCodeCache(freeCache)
	crepo := repository.NewCacheCodeRepository(cc)
	return service.NewCodeCacheService(crepo, initMemorySms())
}

func initMemorySms() sms.Service {
	return localsms.NewService()
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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
		// 允许前端访问后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
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

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//// 引入限流插件，当前限流规则，100QPS
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	//useSession(server)
	useJWT(server)

	return server
}

func useJWT(server *gin.Engine) {
	login := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	// 存储数据的，也就是userId存储在哪里
	// 直接存cookie
	store := cookie.NewStore([]byte("secret"))
	// 基于内存的实现
	//store := memstore.NewStore([]byte("jBxoQWRS5L9vYr$mYq5U9d5BRPHfSBAe"), []byte("O@Gpunh7SPVuLYT^WYBaxDjFUep4THgM"))
	// 基于redis实现:
	//		第一个参数是最大空闲连接数
	//		第二个参数是tcp，不太可能是udp
	//		第三、四个参数是连接信息和密码
	//		第五、六个参数是两个key
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("jBxoQWRS5L9vYr$mYq5U9d5BRPHfSBAe"),
	//	[]byte("O@Gpunh7SPVuLYT^WYBaxDjFUep4THgM"))
	//if err != nil {
	//	panic(err)
	//}
	// 两个middleware,一个用来初始化session，一个用来登录校验
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
