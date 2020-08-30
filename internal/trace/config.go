package trace

import (
	"fmt"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	ot "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"io"
	"time"
)

type holder struct {
	closer io.Closer
	tracer ot.Tracer
}

var (
	httpTimeout = 5 * time.Second
	poolSpans   = jaeger.TracerOptions.PoolSpans(false)
	logger      = spanLogger{}
)

func Configure(serviceName string, options *Options) (io.Closer, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	reporters := make([]jaeger.Reporter, 0, 3)
	sampler, err := jaeger.NewProbabilisticSampler(options.SamplingRate)
	if err != nil {
		return nil, fmt.Errorf("could not build trace sampler: %v", err)
	}
	if options.JaegerURL != "" {
		reporters = append(reporters, jaeger.NewRemoteReporter(transport.NewHTTPTransport(options.JaegerURL, transport.HTTPTimeout(httpTimeout))))
	}
	if options.LogTraceSpans {
		reporters = append(reporters, logger)
	}
	var rep jaeger.Reporter
	switch len(reporters) {
	case 0:
		return holder{}, nil
	case 1:
		rep = reporters[0]
	default:
		rep = jaeger.NewCompositeReporter(reporters...)
	}
	var tracer ot.Tracer
	var closer io.Closer
	tracer, closer = jaeger.NewTracer(serviceName, sampler, rep, poolSpans, jaeger.TracerOptions.Gen128Bit(true))
	//  NOTE: global side effect!
	ot.SetGlobalTracer(tracer)

	return holder{
		closer: closer,
		tracer: tracer,
	}, nil

}

func (h holder) Close() error {
	if ot.GlobalTracer() == h.tracer {
		ot.SetGlobalTracer(ot.NoopTracer{})
	}
	var err error
	if h.closer != nil {
		err = h.closer.Close()
	}
	return err
}

type spanLogger struct{}

func (spanLogger) Report(span *jaeger.Span) {
	logx.Info("Reporting span operation:%s,span:%s", span.OperationName(), span.String())
}
func (spanLogger) Close() {
}
func (spanLogger) Error(msg string) {
	logx.Error(msg)
}
func (spanLogger) Info(msg string, args ...interface{}) {
	logx.Info(msg, args...)
}
