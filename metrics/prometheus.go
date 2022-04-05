package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	LatestValue.Set(-1)
	prometheus.MustRegister(LatestValue)
	prometheus.MustRegister(TotalRequest)
	prometheus.MustRegister(LatestTime)
	prometheus.MustRegister(EndpointAvgResponseTime)
	prometheus.MustRegister(EndpointResponseTime)
	prometheus.MustRegister(EndpointResponseTimeHistogram)
	prometheus.MustRegister(HackCreateUserOnFollow)
	prometheus.MustRegister(HackCreateFollowOnUnfollow)
	prometheus.MustRegister(HackCreateUserOnAddMessage)
}
