package cinit

import (
	"github.com/janiokq/Useless-blog/internal/metrics"
	"time"
)

func metricsInit(sn string) {
	if Config.Metrics.Enable == "yes" {
		//  Push模式
		m := metrics.NewMetrics()
		//  e.Use(api.MetricsFunc(m))
		m.MemStats()
		m.InfluxDBWithTags(
			time.Duration(Config.Metrics.Duration)*time.Second,
			Config.Metrics.URL,
			Config.Metrics.Database,
			Config.Metrics.UserName,
			Config.Metrics.Password,
			Config.Metrics.Measurement,
			map[string]string{"service": sn},
		)

	}
}
