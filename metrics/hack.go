package metrics

import "github.com/prometheus/client_golang/prometheus"

var HackCreateUserOnFollow = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "minitwit_total_hack_create_user_on_follow",
		Help: "total number of times that a user have been created on follow",
	},
)

var HackCreateFollowOnUnfollow = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "minitwit_total_hack_create_follow_on_unfollow",
		Help: "total number of times that a follow entry has been created on unfollow",
	},
)


var HackCreateUserOnAddMessage = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "minitwit_total_hack_create_user_on_add_message",
		Help: "total number of times that a user have been created on add message",
	},
)


