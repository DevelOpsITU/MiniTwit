package metrics

import "github.com/prometheus/client_golang/prometheus"

var EndpointAvgResponseTime = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "minitwit_endpoint_avg_responsetime",
		Help: "The response times of endpoints",
	},
	[]string{"code", "method", "url"},
)

var EndpointResponseTime = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "minitwit_endpoint_responsetime",
		Help: "The response times of endpoints",
	},
	[]string{"code", "method", "url"},
)

var TotalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "minitwit_total_http",
		Help: "total number of http requests",
	},
	[]string{"code", "method", "url"},
)
