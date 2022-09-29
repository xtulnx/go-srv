package config

// GITTAG ?= $(shell git log --pretty=format:"%h:%cd" --date=format:"%Y-%m-%dT%H:%M:%S" -1)
// BUILD_TIME ?= `date +%FT%T%z`
// BUILD_VERSION ?= 0.1.0.1
// APP ?= go-srv
// ConfigName ?= ${APP}
// AppNote ?= "中间层"
//
//
// ENV_PKG=github.com/xtulnx/go-srv/appinfo
//
// ldflagsWeb+=\
//	-X '${ENV_PKG}.IsDebug=0' \
//	-X '${ENV_PKG}.GitTag=${GITTAG}' \
//	-X '${ENV_PKG}.Version=${BUILD_VERSION}' \
//	-X '${ENV_PKG}.BuildTime=${BUILD_TIME}' \
//	-X '${ENV_PKG}.AppName=${APP}' \
//	-X '${ENV_PKG}.AppNote=${AppNote}' \
//	-X '${ENV_PKG}.ConfigName=${ConfigName}' \
//	-w -s

// 版本信息
var (
	//IsDebug = "0"
	IsDebug   = "1"                        // 调试模式
	Version   = "0.1.0"                    // 版本号
	GitTag    = "2022.09.26.debug"         // 代码版本
	BuildTime = "2022-09-26T15:51:57+0800" // 编译时间

	ConfigName = "go-srv" // 配置文件名
	AppName    = "go-srv" // 应用名称
	AppNote    = "中间层"    // 应用简单描述
)
