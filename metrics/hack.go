package metrics

import "github.com/prometheus/client_golang/prometheus"

var HackCreateUserOnFollow = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "minitwit_total_hack_create_user_on_follow",
		Help: "total number of times that a user have been created on follow",
	},
)
