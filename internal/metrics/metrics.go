package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Counters

	// RequestsTotal (by endpoint)
	RequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total requests count",
	}, []string{"endpoint"})

	// RequestsErrors (by endpoint and type)
	RequestsErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_errors",
		Help: "Requests errors count",
	}, []string{"endpoint", "type"})

	// Histograms

	// RequestDurationTotal (by vertical)
	RequestDurationTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "request_duration_seconds",
		Help: "Total request duration count",
	}, []string{"endpoint"})
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestsErrors)
	prometheus.MustRegister(RequestDurationTotal)
}
