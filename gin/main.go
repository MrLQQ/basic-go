package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.Use(func(context *gin.Context) {
		println("这是第一个middleware")
	}, func(context *gin.Context) {
		println("这是第二个middleware")
	})
	// 静态路由
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})

	server.POST("/login", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,login")
	})

	// 参数路由
	// 路由参数
	server.GET("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "hello,"+name)
	})
	//查询参数
	server.GET("/order", func(ctx *gin.Context) {
		id := ctx.Query("id")
		ctx.String(http.StatusOK, "订单ID为,"+id)
	})

	// 通配符路由
	server.GET("/views/*.html", func(ctx *gin.Context) {
		view := ctx.Param(".html")
		ctx.String(http.StatusOK, "您正在浏览"+view)
	})

	// 如果不传参数，会监听8080端口
	server.Run(":8081")
}
