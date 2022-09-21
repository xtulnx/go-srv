// Package wechatkit 微信相关，依赖:
//
//	github.com/silenceper/wechat/v2
package wechatkit

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

type WechatKit struct {
	cfg         WechatConfig
	wc          *wechat.Wechat
	miniprogram *miniprogram.MiniProgram
}

func NewWechat(cfg WechatConfig) (*WechatKit, error) {
	wc := wechat.NewWechat()
	wc.SetCache(cache.NewMemory())
	_miniprogram := wc.GetMiniProgram(&miniConfig.Config{AppID: cfg.AppID, AppSecret: cfg.AppSecret})
	return &WechatKit{
		cfg:         cfg,
		wc:          wc,
		miniprogram: _miniprogram,
	}, nil
}

func (J *WechatKit) Wechat() *wechat.Wechat {
	return J.wc
}

func (J *WechatKit) Cfg() *WechatConfig {
	return &J.cfg
}

func (J *WechatKit) MiniProgram() *miniprogram.MiniProgram {
	return J.miniprogram
}
