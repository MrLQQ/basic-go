package ioc

import (
	"basic-go/webook/internal/web"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/ginx/middleware/ratelimit"
	"basic-go/webook/pkg/limiter"
	"basic-go/webook/pkg/logger"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	wechatHdl *web.OAuth2WechatHandler,
	artHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable,
	hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHandler(),
		func(ctx *gin.Context) {
			println("这是我的middleware……")
		},
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)).Build(),
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
			l.Debug("", logger.Field{Key: "req", Value: al})
		}).AllowReqBody().AllowRespBody().Build(),
		// 使用JWT
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
	}
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowAllOrigins: true,
		//允许的请求源
		AllowOrigins: []string{"http://localhost:3000"},
		//是否允许带上用户认证信息（比如cookie）
		AllowCredentials: true,
		//业务请求中可以带上的头
		AllowHeaders: []string{"Content-Type", "authorization"},
		// 允许前端访问后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
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
	})
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
