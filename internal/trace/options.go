package trace

import "errors"

type Options struct {
	JaegerURL     string
	LogTraceSpans bool
	SamplingRate  float64
}

func (o *Options) Validate() error {
	if o.JaegerURL == "" {
		return errors.New("can't have Jaeger outputs active simultaneously")
	}
	if o.SamplingRate > 1.0 || o.SamplingRate < 0.0 {
		return errors.New("sampling rate must be in the range: [0.0, 1.0]")
	}
	return nil
}

func (o *Options) TracingEnabled() bool {
	return o.JaegerURL != "" || o.LogTraceSpans
}
