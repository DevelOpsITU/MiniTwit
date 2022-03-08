package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prometheusHandlers(router *gin.Engine) {
	LatestValue.Set(69)
	router.GET("/metrics", prometheusHandler())
}

func init() {
	prometheus.MustRegister(LatestValue)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// TODO: This is only for testing setup - remove later
var LatestValue = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "minitwit_latest",
})
