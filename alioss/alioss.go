package alioss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"path/filepath"
	"strings"
)

type AliOss struct {
	cfg AliOssConfig

	client *oss.Client
}

func NewClient(cfg AliOssConfig) (*AliOss, error) {
	c := AliOss{cfg: cfg}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (A *AliOss) init() error {
	if A.cfg.AccessKeyId == "" || A.cfg.AccessKeySecret == "" {
		return fmt.Errorf("无效 oss 授权")
	}
	client, err := oss.New(A.cfg.EndPoint, A.cfg.AccessKeyId, A.cfg.AccessKeySecret)
	if err != nil {
		return err
	}
	A.client = client
	return nil
}
func (A *AliOss) GetClient() (*oss.Client, error) {
	return A.client, nil
}

// UploadKey 上传后的文件对象路径
func (A *AliOss) UploadKey(subKey string) string {
	objectName := filepath.Join(A.cfg.UploadDir, subKey)
	return objectName
}

// UploadFile 上传文件到 OSS
func (A *AliOss) UploadFile(localFile string, objectName string) error {
	bucket, err := A.client.Bucket(A.cfg.Bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObjectFromFile(objectName, localFile)
	return err
}

// GenUrl 完整的对外可访问的URL
func (A *AliOss) GenUrl(key string) string {
	if A.cfg.Host == "" {
		return key
	}
	if key == "" {
		return ""
	}
	if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	return PathJoin(A.cfg.Host, key)
}

func PathJoin(a, b string) string {
	if a == "" {
		return b
	} else if b == "" {
		return a
	}
	cH := a[len(a)-1]
	cK := b[0]
	if cH == '/' && cK == '/' {
		return a + b[1:]
	}
	if cH != '/' && cK != '/' {
		return a + "/" + b
	}
	return a + b
}

// SplitKey 尝试分离对象路径
func (A *AliOss) SplitKey(u string) (h, k string) {
	if strings.HasPrefix(u, A.cfg.Host) {
		return A.cfg.Host, u[len(A.cfg.Host)+1:]
	}
	return "", k
}

// GenUrlThumbnail 缩略图
func (A *AliOss) GenUrlThumbnail(key string) string {
	u := A.GenUrl(key)
	if u == "" {
		return u
	}
	return u + "?x-oss-process=image/auto-orient,0/quality,Q_76/resize,h_192,w_192"
}
