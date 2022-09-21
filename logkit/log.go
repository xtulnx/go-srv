// Package logkit 基于 logrus 的封装:
//
//	日志按天拆分
package logkit

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

// LogConfig 日志的配置
type LogConfig struct {
	// 日志等级: panic fatal error warn info debug trace
	Level string `json:"level" toml:"level"`

	// 日志格式: json text(默认)
	Format string `json:"format" toml:"format"`

	// 输出到文件（如果dir为空，则不输出到文件）
	LogDir string `json:"logdir" toml:"logdir"`

	// 所有日志的输出文件名
	LogFileAll string `json:"logfileall" toml:"logfileall"`

	// 获取日志的输出文件名
	LogFileWarn string `json:"logfilewarn" toml:"logfilewarn"`
}

var (
	std = logrus.StandardLogger()
)

type Fields = logrus.Fields

func init() {
	std.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	std.SetLevel(logrus.InfoLevel)
	//std.SetOutput(os.Stdout)
}

// Init 初始化
func Init(conf LogConfig) error {
	if conf.Level != "" {
		if level, err := logrus.ParseLevel(conf.Level); err == nil {
			std.SetLevel(level)
		} else {
			return err
		}
	}
	if conf.Format == "json" {
		std.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05",
		})
	}
	// 输出
	if conf.LogDir != "" {
		if conf.LogFileAll != "" {
			AddLogFile(std, conf.LogDir, conf.LogFileAll, time.Hour*24*30, time.Hour*24, LogLevels(logrus.DebugLevel))
		}
		if conf.LogFileWarn != "" {
			AddLogFile(std, conf.LogDir, conf.LogFileWarn, time.Hour*24*30, time.Hour*24, LogLevels(logrus.WarnLevel))
		}
	}
	return nil
}

func LogLevels(baseLevel logrus.Level) []logrus.Level {
	if int(baseLevel) < len(logrus.AllLevels) {
		l1 := make([]logrus.Level, int(baseLevel)+1)
		copy(l1, logrus.AllLevels[0:int(baseLevel)+1])
		return l1
	} else {
		return []logrus.Level{}
	}
}

// AddLogFile 关联日志文件. 默认使用 TextFormatter
//
//	_log: 日志，用于 .AddHook
//	logPath: 目录
//	logFileName: 日志文件名
//	maxAge: 最长保存时间
//	rotationTime: 分割时间
//	levels: 关注的日志等级
func AddLogFile(_log *logrus.Logger, logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration, levels []logrus.Level) {
	if len(levels) == 0 {
		fmt.Printf("空的日志等级")
		return
	}
	baseLogPath := path.Join(logPath, logFileName)

	fmt.Printf("日志目录: %s, 等级：%v\n", baseLogPath, levels)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logFileName),
		// WithRotationTime设置日志分割的时间
		rotatelogs.WithRotationTime(rotationTime),

		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		rotatelogs.WithMaxAge(maxAge),
		//rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		fmt.Printf("config local file system for logger error: %v", errors.WithStack(err))
	}

	wm := make(lfshook.WriterMap)
	for _, l := range levels {
		wm[l] = writer
	}
	lfsHook := lfshook.NewHook(wm, &logrus.TextFormatter{DisableColors: true})
	_log.AddHook(lfsHook)
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// ForTask 添加 task 字段
func ForTask(name string) *logrus.Entry {
	return std.WithField("task", name)
}
