package controllers

import (
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
}

// HandleRESTRequests - handles the rest requests
func HandleRESTRequests() {

	router := gin.Default()
	router.SetTrustedProxies(nil)

	for _, handler := range HttpHandlers {
		handler.(func(engine *gin.Engine))(router)
	}

	router.Run(":8080")

}
