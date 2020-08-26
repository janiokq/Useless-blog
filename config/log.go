package config

import (
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/janiokq/Useless-blog/internal/utils/logx/configlog"
	"github.com/janiokq/Useless-blog/internal/utils/logx/plugins/zaplog"
)

// 初始化日志,可以再这里初始化不同日志引擎的日志 、、 zap logrous// 初始化zap// 设置日志引擎为刚初始化的
func logInit() {
	logx.SetLogger(zaplog.New(
		configlog.WithProjectName(Config.Service.Name),
		configlog.WithLogPath(Config.Log.Path),
		configlog.WithLogName(Config.Service.Name),
		configlog.WithMaxAge(Config.Log.MaxAge),
		configlog.WithMaxSize(Config.Log.MaxSize),
		configlog.WithIsStdout(Config.Log.IsStdout),
		configlog.WithLogLevel(Config.Log.LogLevel),
	))
}
