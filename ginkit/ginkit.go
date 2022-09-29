// Package ginkit 封装一些通用方法
//
//	github.com/gin-gonic/gin
//
// ../ginjson
//
// json 有处理，需要编译参数 -tags jsoniter
package ginkit

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(g *gin.Engine) {

	//中间件
	middlewares := []gin.HandlerFunc{
		gin.Logger(),
		gin.Recovery(),
		NoCache,
		Options,
		Secure,
	}
	g.Use(middlewares...)

	g.RemoveExtraSlash = true

	//404处理
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "该路径不存在")
	})

	//g.LoadHTMLGlob("template/**/*") //加载模板路径

	//健康检查中间件
	//g.GET("/", service.Index)               //主页
	//g.GET("/gooss", service.Gooss)          //oss信息
	//g.GET("/ossupload", service.OssUpload)  //上传oss
	//g.POST("/ossupload", service.OssUpload) //上传oss
}
