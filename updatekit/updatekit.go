// Package updatekit 热更新封装
//
//	依赖 overseer
package updatekit

import (
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"time"
)

type UpdateConfig struct {
	Enabled bool

	Address   string
	Addresses []string

	Fetcher  string
	Path     string
	Url      string
	Interval time.Duration
}

func (cfg UpdateConfig) GetFetcher() (f1 fetcher.Interface) {
	switch cfg.Fetcher {
	case "file":
		f1 = &fetcher.File{
			Path:     cfg.Path,
			Interval: cfg.Interval,
		}
	case "http":
		f1 = &fetcher.HTTP{
			URL:          cfg.Url,
			Interval:     cfg.Interval,
			CheckHeaders: nil,
		}
	case "github":
		f1 = &fetcher.Github{
			User:     "",
			Repo:     "",
			Interval: cfg.Interval,
			Asset:    nil,
		}
	case "s3":
		f1 = &fetcher.S3{
			Access:      "",
			Secret:      "",
			Region:      "",
			Bucket:      "",
			Key:         "",
			Interval:    cfg.Interval,
			HeadTimeout: 0,
			GetTimeout:  0,
		}
	default:
	}
	if f1 != nil {

	}
	return f1
}

func Run(cfg UpdateConfig, program func(state overseer.State)) {
	f1 := cfg.GetFetcher()
	c1 := overseer.Config{
		Required:            false,
		Program:             program,
		Address:             cfg.Address,
		Addresses:           cfg.Addresses,
		RestartSignal:       nil,
		TerminateTimeout:    0,
		MinFetchInterval:    0,
		PreUpgrade:          nil,
		Debug:               false,
		NoWarn:              false,
		NoRestart:           false,
		NoRestartAfterFetch: false,
		Fetcher:             f1,
	}
	overseer.Run(c1)
}
