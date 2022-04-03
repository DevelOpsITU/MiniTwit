package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var PostgresData = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "postgres_data",
	}, []string{"table"},
)
