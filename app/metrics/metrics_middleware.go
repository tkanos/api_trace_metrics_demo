package metrics

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	requestFailedName = "failed_count"
	latencyName       = "duration_milliseconds"
)

type metricsService struct {
	requestFailed *prometheus.CounterVec
	latency       *prometheus.HistogramVec
	serviceName   string
}

// NewMetricsMiddleware creates a layer of service that add metrics capability
func NewMetricsMiddleware(serviceName string, next endpoint.Endpoint) endpoint.Endpoint {
	m := metricsMiddleware(serviceName)
	return m.instrumentation(next)
}

func metricsMiddleware(name string) *metricsService {
	var m metricsService
	fieldKeys := []string{"method"}

	m.requestFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "request",
			Name:      requestFailedName,
			Help:      "Number of requests failed",
		}, fieldKeys)
	prometheus.MustRegister(m.requestFailed)

	m.latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "request",
			Name:      latencyName,
			Help:      "Total duration in miliseconds.",
			Buckets:   []float64{50, 100, 200, 1000},
			//ConstLabels: prometheus.Labels{"":"","":""},
		}, fieldKeys)
	prometheus.MustRegister(m.latency)

	m.serviceName = name

	return &m
}

func (m *metricsService) instrumentation(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		start := time.Now()
		// add metrics to this method
		defer func(start time.Time) {
			m.latency.WithLabelValues(m.serviceName).Observe(time.Since(start).Seconds() * 1e3)
		}(start)
		// If error is not empty, we add to metrics that it failed
		response, err = next(ctx, request)
		if err != nil {
			m.requestFailed.WithLabelValues(m.serviceName).Inc()
		}
		return
	}
}
