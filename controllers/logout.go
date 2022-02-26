package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minitwit/models"
	"net/http"
)

func logoutHandlers(router *gin.Engine) {

	router.GET("/logout", func(c *gin.Context) {
		handleLogout(c.Writer, c.Request, c)
	})
}

func handleLogout(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	// reset cookie
	g := models.Session{
		User:     models.User{},
		Message:  true,
		Messages: []string{"You were logged out"},
	}
	data, _ := json.Marshal(g)
	c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}
