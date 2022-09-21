package alioss

type AliOssConfig struct {
	AccessKeyId     string
	AccessKeySecret string //
	EndPoint        string // 可以使用云端内网域名
	Bucket          string //
	UploadDir       string // 上传目录，如 j00/demo/
	Host            string // 主机名称，格式为 bucketname.endpoint ，如 https://j00-demo.oss-cn-shenzhen.aliyuncs.com
	Callback        string // 服务端授权上传时的回调地址（完整地址）

	HostAlias []string // 额外名称
}
