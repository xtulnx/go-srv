package alioss

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/xtulnx/go-srv/errno"
	"github.com/xtulnx/go-srv/utils"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// https://help.aliyun.com/document_detail/31837.html 访问域名和数据中心
// https://help.aliyun.com/document_detail/31926.htm 服务端签名后直传
// https://help.aliyun.com/document_detail/31988.htm 关于Object操作/基础操作/PostObject
//
// https://help.aliyun.com/document_detail/31989.html 关于Object操作/基础操作/Callback

// https://github.com/aliyun/alibaba-cloud-sdk-go Alibaba Cloud SDK for Go 让您不用复杂编程即可访问云服务器、云监控等多个阿里云服务

// https://juejin.cn/post/7082240703438258206 golang对接阿里云私有Bucket上传图片、授权访问图片
//
//  https://help.aliyun.com/document_detail/91818.html
//   首页 >对象存储 OSS >最佳实践 >网站与移动应用 >Web端上传数据至OSS >服务端签名直传并设置上传回调 >Go
//
//  https://help.aliyun.com/document_detail/31837.html 访问域名（Endpoint）/访问域名和数据中心
//

// 服务端授权方式

// PolicyToken 授权策略 token
type PolicyToken struct {
	AccessKeyId string `json:"accessid"`  // 用户请求的AccessKey ID。授权 ID
	Signature   string `json:"signature"` // 对Policy签名后的字符串。
	Policy      string `json:"policy"`    // 用户表单上传的策略（Policy），Policy为经过Base64编码过的字符串。
	Callback    string `json:"callback"`  // 文件上传结果回调
	Host        string `json:"host"`      // 用户发送上传请求的域名。oss 的 endpoint
	Expire      int64  `json:"expire"`    // 由服务器端指定的Policy过期时间，格式为Unix时间戳（自UTC时间1970年01月01号开始的秒数）。

	Directory string `json:"dir,omitempty"` // 限制上传的文件前缀。直接拼接，**不要再加「/」**
	FileKey   string `json:"key,omitempty"` // 指定文件路径。如果有值，则相当于本次token只能用于上传一个文件
}

type ConfigStruct struct {
	// Expiration项指定了Policy的过期时间，以ISO8601 GMT时间表示。
	// 例如2014-12-01T12:00:00.000Z指定了Post请求必须在2014年12月1日12点之前。
	Expiration string `json:"expiration"`
	// Conditions是一个列表，可以用于指定Post请求的表单域的合法值。
	//  https://help.aliyun.com/document_detail/31988.htm
	Conditions [][]interface{} `json:"conditions"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

// GetPolicyToken 获取 上传授权 token
//
//	prefix 路径前缀，在整体的配置 UploadDir 之下，如果是 以 「/」结尾，则表示"目录"
func (A *AliOss) GetPolicyToken(prefix, callbackBody string) PolicyToken {
	uploadDir := filepath.Join(A.cfg.UploadDir, prefix)
	callbackUrl := A.cfg.Callback
	expire_end := time.Now().Unix() + 30

	//create post policy json
	config := ConfigStruct{
		Expiration: time.Unix(expire_end, 0).UTC().Format("2006-01-02T15:04:05Z"),
		Conditions: [][]interface{}{
			// 指定前缀
			{"starts-with", "$key", uploadDir},
			// 指定路径的方式
			//{"eq", "$key", ""},
			// 限制上传文件大小。
			{"content-length-range", 0, 30 * 1024 * 1024},
		},
	}

	//calucate signature
	result, _ := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)

	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(A.cfg.AccessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callbackUrl
	callbackParam.CallbackBody = strings.Join([]string{
		"bucket=${bucket}",
		"filename=${object}",
		"size=${size}",
		"mimeType=${mimeType}",
		"etag=${etag}",

		// imageInfo针对图片格式，如果为非图片格式，这些为空
		"height=${imageInfo.height}",
		"width=${imageInfo.width}",
		"format=${imageInfo.format}",
	}, "&")
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	if callbackBody != "" {
		callbackParam.CallbackBody += "&" + callbackBody
	}
	callback_str, _ := json.Marshal(callbackParam)
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	policyToken := PolicyToken{
		AccessKeyId: A.cfg.AccessKeyId,
		Host:        A.cfg.Host,
		Expire:      expire_end,
		Signature:   signedStr,
		Directory:   uploadDir,
		Policy:      debyte,
		Callback:    callbackBase64,
	}
	return policyToken
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// getPublicKey : Get PublicKey bytes from Request.URL
func getPublicKey(logger utils.Logger, r *http.Request) ([]byte, error) {
	var bytePublicKey []byte
	if publicKeyURLBase64 := r.Header.Get("x-oss-pub-key-url"); publicKeyURLBase64 == "" {
		logger.Warn("GetPublicKey from Request header failed :  No x-oss-pub-key-url field. ")
		return bytePublicKey, errors.New("no x-oss-pub-key-url field in Request header ")
	} else if publicKeyURL, e1 := base64.StdEncoding.DecodeString(publicKeyURLBase64); e1 != nil {
		return nil, e1
		// TODO: 考虑使用缓存
	} else if responsePublicKeyURL, e2 := http.Get(string(publicKeyURL)); e2 != nil {
		logger.Warnf("Get PublicKey Content from URL failed : %s \n", e2.Error())
		return bytePublicKey, e2
	} else {
		bytePublicKey, e2 = ioutil.ReadAll(responsePublicKeyURL.Body)
		_ = responsePublicKeyURL.Body.Close()
		if e2 != nil {
			logger.Warnf("Read PublicKey Content from URL failed : %s \n", e2.Error())
			return bytePublicKey, e2
		}
	}
	return bytePublicKey, nil
}

// getAuthorization : decode from Base64String
func getAuthorization(logger utils.Logger, r *http.Request) ([]byte, error) {
	strAuthorizationBase64 := r.Header.Get("authorization")
	if strAuthorizationBase64 == "" {
		logger.Warn("Failed to get authorization field from request header. ")
		return nil, errors.New("no authorization field in Request header")
	}
	byteAuthorization, e1 := base64.StdEncoding.DecodeString(strAuthorizationBase64)
	if e1 != nil {
		return nil, e1
	}
	return byteAuthorization, nil
}

// getMD5FromNewAuthString : Get MD5 bytes from Newly Constructed Authrization String.
//
//	Construct the New Auth String from URI+Query+Body
func getMD5FromNewAuthString(logger utils.Logger, r *http.Request) (byteMD5 []byte, params url.Values, err error) {
	var u1 *url.URL
	if uri0 := r.Header.Get("x-real-uri"); uri0 != "" {
		u1, _ = url.Parse(uri0)
	}
	if u1 == nil {
		u1 = r.URL
	}

	bodyContent, err := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		logger.Warnf("Read Request Body failed : %s \n", err.Error())
		return byteMD5, nil, err
	}
	strCallbackBody := string(bodyContent)

	// Generate New Auth String prepare for MD5
	strAuth := ""
	if r.URL.RawQuery == "" {
		strAuth = strings.Join([]string{u1.Path, "\n", strCallbackBody}, "")
	} else {
		strAuth = strings.Join([]string{u1.Path, "?", r.URL.RawQuery, "\n", strCallbackBody}, "")
	}

	params = url.Values{}
	if r.URL.RawQuery != "" {
		if v1, e1 := url.ParseQuery(r.URL.RawQuery); e1 != nil {
		} else {
			for a := range v1 {
				if b := v1.Get(a); b != "" {
					params.Add(a, b)
				}
			}
		}
	}
	if strCallbackBody != "" {
		if v1, e1 := url.ParseQuery(strCallbackBody); e1 != nil {
		} else {
			for a := range v1 {
				if b := v1.Get(a); b != "" {
					params.Add(a, b)
				}
			}
		}
	}

	// Generate MD5 from the New Auth String
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(strAuth))
	byteMD5 = md5Ctx.Sum(nil)
	return
}

/*  VerifySignature
*   VerifySignature需要三个重要的数据信息来进行签名验证： 1>获取公钥PublicKey;  2>生成新的MD5鉴权串;  3>解码Request携带的鉴权串;
*   1>获取公钥PublicKey : 从RequestHeader的"x-oss-pub-key-url"字段中获取 URL, 读取URL链接的包含的公钥内容， 进行解码解析， 将其作为rsa.VerifyPKCS1v15的入参。
*   2>生成新的MD5鉴权串 : 把Request中的url中的path部分进行urldecode， 加上url的query部分， 再加上body， 组合之后进行MD5编码， 得到MD5鉴权字节串。
*   3>解码Request携带的鉴权串 ： 获取RequestHeader的"authorization"字段， 对其进行Base64解码，作为签名验证的鉴权对比串。
*   rsa.VerifyPKCS1v15进行签名验证，返回验证结果。
* */
func verifySignature(logger utils.Logger, bytePublicKey []byte, byteMd5 []byte, authorization []byte) bool {
	pubBlock, _ := pem.Decode(bytePublicKey)
	if pubBlock == nil {
		logger.Warn("Failed to parse PEM block containing the public key")
		return false
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if (pubInterface == nil) || (err != nil) {
		logger.Warnf("x509.ParsePKIXPublicKey(publicKey) failed : %s \n", err.Error())
		return false
	}
	pub := pubInterface.(*rsa.PublicKey)

	errorVerifyPKCS1v15 := rsa.VerifyPKCS1v15(pub, crypto.MD5, byteMd5, authorization)
	if errorVerifyPKCS1v15 != nil {
		logger.Warnf("\nSignature Verification is Failed : %s \n", errorVerifyPKCS1v15.Error())
		//printByteArray(byteMd5, "AuthMd5(fromNewAuthString)")
		//printByteArray(bytePublicKey, "PublicKeyBase64")
		//printByteArray(authorization, "AuthorizationFromRequest")
		return false
	}

	logger.Warn("\nSignature Verification is Successful. \n")
	return true
}

// VerifyCallback 回调校验
func (A *AliOss) VerifyCallback(logger utils.Logger, r *http.Request) (url.Values, error) {
	bytePublicKey, err := getPublicKey(logger, r)
	if err != nil {
		return nil, err
	}
	byteAuthorization, err := getAuthorization(logger, r)
	if err != nil {
		return nil, err
	}
	// Get MD5 bytes from Newly Constructed Authrization String.
	byteMD5, body, err := getMD5FromNewAuthString(logger, r)
	if err != nil {
		return nil, err
	}

	// verifySignature and response to client
	if verifySignature(logger, bytePublicKey, byteMD5, byteAuthorization) {
		// do something you want accoding to callback_body ...

		return body, nil
	}
	return nil, errno.BadRequest.SetMsg("签名校验不通过")
}

/*
const host = '<host>';
const signature = '<signatureString>';
const ossAccessKeyId = '<accessKey>';
const policy = '<policyBase64Str>';
const key = '<object name>';
const securityToken = '<x-oss-security-token>';
const filePath = '<filePath>'; // 待上传文件的文件路径。
wx.uploadFile({
  url: host, // 开发者服务器的URL。
  filePath: filePath,
  name: 'file', // 必须填file。
  formData: {
    key,
    policy,
    OSSAccessKeyId: ossAccessKeyId,
    signature,
    // 'x-oss-security-token': securityToken // 使用STS签名时必传。
  },
  success: (res) => {
    if (res.statusCode === 204) {
      console.log('上传成功');
    }
  },
  fail: err => {
    console.log(err);
  }
});
*/
