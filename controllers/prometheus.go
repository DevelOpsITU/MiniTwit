package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/********************************************************
*														*
*					Prometheus metrics					*
*														*
********************************************************/
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

/************** REMEMBER TO REGISTER *******************/
func init() {
	LatestValue.Set(-1)
	prometheus.MustRegister(LatestValue)
	prometheus.MustRegister(TotalRequest)
	prometheus.MustRegister(LatestTime)
	prometheus.MustRegister(EndpointAvgResponseTime)
	prometheus.MustRegister(EndpointResponseTime)
}

/********************************************************
*														*
*					Controller Handler					*
*														*
********************************************************/
func prometheusHandlers(router *gin.Engine) {
	router.GET("/metrics", prometheusHandler())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

/********************************************************
*														*
*				Gin Context Middleware					*
*														*
********************************************************/
// error handling middleware (could be reused for logging as well)
func HttpGinMiddleware(c *gin.Context) {

	// get request from context
	request := c.Request

	// maybe sort away /metrics
	// ------------------------

	// measure time of request
	start := time.Now()
	c.Next()
	handleTime := time.Since(start)

	// get status code from gin context
	statuscode := fmt.Sprint(c.Writer.Status())

	// check if failed
	for _, err := range c.Errors {
		// log error (change to use actual logger)
		fmt.Println(err.Error())
		// TotalRequest.WithLabelValues(statuscode, request.Method).Inc()
	}

	// Replaces parameters with :key so that they are grouped under one
	url := request.URL.Path
	for _, p := range c.Params {
		url = strings.ReplaceAll(url, p.Value, ":"+p.Key)
	}

	EndpointAvgResponseTime.WithLabelValues(statuscode, request.Method, url).Observe(float64(handleTime.Nanoseconds()))
	EndpointResponseTime.WithLabelValues(statuscode, request.Method, url).Set(float64(handleTime.Nanoseconds()))

	LatestTime.Set(float64(handleTime.Nanoseconds()))
	TotalRequest.WithLabelValues(statuscode, request.Method, url).Inc()
}
