package metrics

import "github.com/prometheus/client_golang/prometheus"

var LatestValue = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "minitwit_latest",
	},
)

var LatestTime = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "minitwit_latest_time",
	},
)
