package controllers

import (
	"fmt"
	"math"
	"minitwit/metrics"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

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

	metrics.EndpointAvgResponseTime.WithLabelValues(statuscode, request.Method, url).Observe(float64(handleTime.Nanoseconds()))
	metrics.EndpointResponseTime.WithLabelValues(statuscode, request.Method, url).Set(float64(handleTime.Nanoseconds()))
	metrics.EndpointResponseTimeHistogram.Observe(float64(handleTime.Nanoseconds()) / (math.Pow(10, 9)))

	metrics.LatestTime.Set(float64(handleTime.Nanoseconds()))
	metrics.TotalRequest.WithLabelValues(statuscode, request.Method, url).Inc()
}
