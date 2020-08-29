package cinit

import (
	"github.com/janiokq/Useless-blog/internal/trace"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/opentracing/opentracing-go"
	"io"
)

var c io.Closer

func tracerInit() {
	// 配置
	c = traceingInit(Config.Trace.Address, Config.Trace.LogTraceSpans, Config.Trace.SamplingRate, Config.Service.Name)
	logx.Infof("初始化traceing:%+v", opentracing.GlobalTracer())

}

func traceingInit(jaegerURL string, logTraceSpans bool, samplingRate float64, servicename string) io.Closer {
	cl, err := trace.Configure(servicename, &trace.Options{
		JaegerURL:     jaegerURL,
		LogTraceSpans: logTraceSpans,
		SamplingRate:  samplingRate,
	})
	if err != nil {
		logx.Error(err.Error())
	}
	return cl
}

func tracerClose() {

}
