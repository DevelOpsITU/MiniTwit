package controllers

import (
	"fmt"
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
var LatestValue = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "minitwit_latest",
})

var LatestTime = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "minitwit_latest_time",
	},
)

var TotalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "minitwit_total_http",
		Help: "total number of GET requests",
	},
	[]string{"code", "method"})

/************** REMEMBER TO REGISTER *******************/
func init() {
	prometheus.MustRegister(LatestValue)
	prometheus.MustRegister(TotalRequest)
	prometheus.MustRegister(LatestTime)
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
func ErrorHandler(c *gin.Context) {

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

	LatestTime.Set(float64(handleTime.Nanoseconds()))
	TotalRequest.WithLabelValues(statuscode, request.Method).Inc()
	LatestValue.Add(1)
}
