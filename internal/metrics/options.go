package metrics

import "github.com/rcrowley/go-metrics"

type Option func(c *Options)

type Options struct {
	Registry metrics.Registry
	Prefix   string
}

func applyOptions(options ...Option) Options {
	opts := Options{
		Registry: metrics.DefaultRegistry,
	}

	for _, option := range options {
		option(&opts)
	}

	return opts
}
