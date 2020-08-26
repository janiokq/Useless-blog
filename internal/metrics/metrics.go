package metrics

import (
	"github.com/rcrowley/go-metrics"
	influxdb "github.com/vrischmann/go-metrics-influxdb"
	"net"
	"time"
)

type Metrics struct {
	opts Options
}

func NewMetrics(options ...Option) *Metrics {
	opts := applyOptions(options...)
	return &Metrics{opts: opts}
}

func (m *Metrics) WithPrefix(s string) string {
	return m.opts.Prefix + "." + s
}

func (m *Metrics) GetRegistry() metrics.Registry {
	return m.opts.Registry
}

func (m *Metrics) MemStats() {
	metrics.RegisterRuntimeMemStats(m.opts.Registry)
	go metrics.CaptureRuntimeMemStats(m.opts.Registry, 5*time.Second)
}

func (m *Metrics) Log(freq time.Duration, l metrics.Logger) {
	go metrics.Log(m.opts.Registry, freq, l)
}

func (m *Metrics) Graphite(freq time.Duration, prefix string, addr *net.TCPAddr) {
	go metrics.Graphite(m.opts.Registry, freq, prefix, addr)
}

func (m *Metrics) InfluxDB(freq time.Duration, url, database, username, password string) {
	go influxdb.InfluxDB(m.opts.Registry, freq, url, database, "measurement", username, password, false)
}

func (m *Metrics) InfluxDBWithTags(freq time.Duration, url, database, username, password string, measurement string, tags map[string]string) {
	go influxdb.InfluxDBWithTags(m.opts.Registry, freq, url, database, measurement, username, password, tags, false)
}
