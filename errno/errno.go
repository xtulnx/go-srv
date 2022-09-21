package errno

import "fmt"

// 2xx成功
// 3xx重定向
// 4xx客户端错误
// 5xx服务器错误

var (
	OK                  = &Errno{200, "OK"}   //
	BadRequest          = &Errno{400, "参数有误"} //
	Unauthorized        = &Errno{401, "请先登录"}
	Forbidden           = &Errno{403, "权限不足"}
	Conflict            = &Errno{409, "已经存在"}
	TooManyRequests     = &Errno{429, "访问太快，请稍候再试"} //
	InternalServerError = &Errno{500, "服务器错误"}      //
	QueryNotFound       = &Errno{550, "记录并不存在"}
	QueryFailed         = &Errno{551, "查询记录失败"}
	ConvertDataFailed   = &Errno{552, "数据格式有误"}
	FailedUpdate        = &Errno{553, "数据更新失败"}
)

// Errno 错误码
type Errno struct {
	Code    int
	Message string
}

func (e Errno) Error() string {
	return e.Message
}

// Err 带链
type Err struct {
	Errno
	Err error
}

func (err Err) Error() string {
	return fmt.Sprintf("Err - code:%d", err.Code)
}

// DecodeErr 取出错误码与描述（不含错误链）。 默认返回 OK 的信息
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	case Errno:
		return typed.Code, typed.Message
	default:
	}
	return InternalServerError.Code, err.Error()
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// SetMsg 更换错误描述
func (e Errno) SetMsg(msg string) Errno {
	return Errno{e.Code, msg}
}

// WithErr 合并添加错误链，保留当前的错误码与描述
func (e Errno) WithErr(err error) *Err {
	return NewErr(e.Code, e.Message, err)
}

// WithMsg 合并添加错误链，保留当前的错误码，重新定义错误描述
func (e Errno) WithMsg(msg string, err error) *Err {
	return NewErr(e.Code, msg, err)
}

func (e Errno) isInner(err error) bool {
	switch err.(type) {
	case *Err:
		return true
	case Err:
		return true
	case Errno:
		return true
	case *Errno:
		return true
	}
	return false
}

// WithErr2 如果 err 是当前包内的错误，则直接引用；
// 否则创建新对象：合并添加错误链，保留当前的错误码和描述
func (e Errno) WithErr2(err error) error {
	if e.isInner(err) {
		return err
	}
	return NewErr(e.Code, e.Message, err)
}

// WithMsg2 如果 err 是当前包内的错误，则直接引用；
// 否则创建新对象：合并添加错误链，保留当前的错误码，重新定义错误描述
func (e Errno) WithMsg2(msg string, err error) error {
	if e.isInner(err) {
		return err
	}
	return NewErr(e.Code, msg, err)
}

// NewErr 创建一个新的错误对象，包括错误链
func NewErr(code int, message string, err error) *Err {
	return &Err{Errno{code, message}, err}
}

func NewErr2(errno *Errno, message string, err error) *Err {
	return &Err{Errno{errno.Code, errno.Message + "," + message}, err}
}

func New(errno *Errno, err error) *Err {
	return &Err{Errno{errno.Code, errno.Message}, err}
}
