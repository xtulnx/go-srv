package ginkit

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xtulnx/go-srv/errno"
	"github.com/xtulnx/go-srv/ginjson"
	"github.com/xtulnx/go-srv/utils"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendResp1(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message, Data: data})
}

type WithCheckParam interface {
	// CheckParam 校验参数
	CheckParam() error
}

type WithBindParam interface {
	// BindParam 读取参数
	BindParam(c *gin.Context) error
}

type WithIPer interface {
	// SetIP 需要设置 IP 地址
	SetIP(string)
}

// SetResponseHeaderForFile 设置输出的文件
func SetResponseHeaderForFile(c *gin.Context, fileName, fileExt string) {
	if fileExt != "" {
		ext := filepath.Ext(fileName)
		if ext != fileExt {
			fileName = fileName + fileExt
		}
	}
	rh := c.Writer.Header()
	mt := mime.TypeByExtension(filepath.Ext(fileName))
	if mt != "" {
		rh.Set("Content-type", mt)
	}
	rh.Set("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
}

// BindParam 解析参数
func BindParam(c *gin.Context, obj interface{}, _log utils.Logger) error {
	// err := c.Bind(obj)
	// 兼容 json 2022.04.20
	b := binding.Default(c.Request.Method, c.ContentType())
	if b == binding.JSON {
		b = ginjson.BindingJSON
	}
	err := c.MustBindWith(obj, b)
	if err != nil {
		_log.Error(err)
		var msg string
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			if e.Struct != "" && e.Field != "" {
				msg = "无效参数" + ": [" + e.Value + "]" + " 字段: " + e.Struct + "." + e.Field + " 类型: " + e.Type.String()
			} else {
				msg = "无效参数" + ": [" + e.Value + "]" + " 类型: " + e.Type.String()
			}
			//} else if e,ok:=err.(*strconv.NumError); ok {
		} else {
			msg = "无效参数" + err.Error()
		}
		return errno.BadRequest.WithMsg(msg, err)
	}

	if iper, ok := obj.(WithIPer); ok {
		iper.SetIP(c.ClientIP())
	}

	// 其他
	if b1, ok := obj.(WithBindParam); ok {
		err = b1.BindParam(c)
	}

	if checker, ok := obj.(WithCheckParam); ok {
		err = checker.CheckParam()
	}

	return nil
}
