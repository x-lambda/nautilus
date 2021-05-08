package log

import (
	"context"
	"time"

	"nautilus/util/conf"
	"nautilus/util/ctxkit"

	"github.com/sirupsen/logrus"
)

// log 全局的 log 对象
var logger *logrus.Logger

// Logger logrus logger封装
type Logger = *logrus.Entry

type Fields = logrus.Fields

// levels 日志等级
var levels = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
}

func init() {
	setLevel()

	logger = logrus.New()

	// TODO 如果有设置log-agent修改logrus.SetOutput()
	if conf.Get("LOG_AGENT") != "" {
	}
}

// setLevel 设置日志等级
func setLevel() {
	level := conf.Get("LOG_LEVEL")
	if level == "" {
		level = "debug"
	}

	logrus.SetLevel(levels[level])
}

// Reset 重置日志等级
func Reset() {
	setLevel()
}

func Get(ctx context.Context) Logger {
	entry := logger.WithFields(logrus.Fields{
		"env":         conf.Env,
		"app_id":      conf.AppID,
		"instance_id": conf.Hostname,
		"trace_id":    ctxkit.GetTraceID(ctx),
		"platform":    ctxkit.GetPlatform(ctx),
		"ip":          ctxkit.GetAccessIP(ctx),
		"ts":          time.Now().Format(time.RFC3339Nano),
	})

	return entry
}
