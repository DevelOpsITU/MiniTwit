package controllers

import (
	"fmt"
	"github.com/gin-contrib/logger"
	"github.com/rs/zerolog"
	"io"
	"minitwit/config"
	"minitwit/log"
	"time"

	"github.com/gin-gonic/gin"
)

var HttpHandlers = []interface{}{
	loginHandlers,
	logoutHandlers,
	userHandlers,
	registerHandlers,
	timelineHandlers,
	staticHandlers,
	addMessageHandlers,
	simulationHandlers,
	prometheusHandlers,
}

// HandleRESTRequests - handles the rest requests
func HandleRESTRequests() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Could not set trusted proxies")
	}

	// TODO: Add proxy header if running in container
	router.Use(logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
			return log.Logger.With().
				Str("path", c.Request.URL.Path).
				Str("code", fmt.Sprint(c.Writer.Status())).
				Dur("latency", latency).
				Logger()
		})))

	for _, handler := range HttpHandlers {
		handler.(func(engine *gin.Engine))(router)
	}

	err = router.Run(fmt.Sprintf(":%s", config.GetConfig().Server.Port))
	if err != nil {
		log.Logger.Error().Err(err).Msg("Could not open the server on " + fmt.Sprintf(":%s", config.GetConfig().Server.Port))
	}

}
