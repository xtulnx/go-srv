package ginjson

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"io"
	"net/http"
)

func init() {
	// 启用 json 兼容处理 ，
	// 还需要在编译参数加上:  -tags jsoniter
	//  如 go run -tags "jsoniter" main.go
	extra.RegisterFuzzyDecoders()
}

func Init() {

}

var (
	jsonLib = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by gin/json package.
	Marshal = jsonLib.Marshal
	// Unmarshal is exported by gin/json package.
	Unmarshal = jsonLib.Unmarshal
	// MarshalIndent is exported by gin/json package.
	MarshalIndent = jsonLib.MarshalIndent
	// NewDecoder is exported by gin/json package.
	NewDecoder = jsonLib.NewDecoder
	// NewEncoder is exported by gin/json package.
	NewEncoder = jsonLib.NewEncoder
)

// BindingJSON 替换Gin默认的binding，支持更丰富JSON功能
var BindingJSON binding.Binding = jsonBinding{}

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return fmt.Errorf("invalid request")
	}
	return decodeJSON(req.Body, obj)
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := NewDecoder(r)
	if binding.EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if binding.EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj interface{}) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
