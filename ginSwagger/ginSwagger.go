package ginSwagger

import (
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
)

func InitRoute(base *gin.RouterGroup, subPath, listen string) {
	apiSwagger := api.Group("/swagger")
	initRouteSwagger(apiSwagger, api.BasePath())

	if _ip, err := net.ResolveTCPAddr("tcp", listen); err == nil {
		var host string
		if _ip.IP == nil {
			host = net.JoinHostPort("127.0.0.1", strconv.Itoa(_ip.Port))
		} else {
			host = _ip.String()
		}
		log.Infof("文档地址： http://%s%s/index.html", host, apiSwagger.BasePath())

	}
}
