package wechatkit

// WechatConfig 微信配置
type WechatConfig struct {
	AppID     string `json:"app_id"`     // appid
	AppSecret string `json:"app_secret"` // appsecret

	MP MiniProgramConfig `json:"mp" toml:"mp" mapstructure:"mp"` // 小程序配置
}

// MiniProgramConfig 小程序相关
type MiniProgramConfig struct {
	// 小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	State string `json:"state" toml:"state" mapstructure:"state"`
	// 订阅配置
	Subscribe SubscribeMpConfig `json:"subscribe" toml:"subscribe" mapstructure:"subscribe"`
}

// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html#method-http

// SubscribeMpConfig 订阅相关
type SubscribeMpConfig struct {
	Templ map[string]TemplSubscribeMpConfig `json:"templ" toml:"templ" mapstructure:"templ"`
}

// TemplSubscribeMpConfig 模板
type TemplSubscribeMpConfig struct {
	// 模板 ID，字符串
	ID string `json:"id" toml:"id" mapstructure:"id"`
	// 字段映射，内部约定的字段 -> 模板中引用的字段名
	Fields map[string]string `json:"fields" toml:"fields" mapstructure:"fields"`
	// 页面跳转
	Page string `json:"page" toml:"page" mapstructure:"page"`
}
