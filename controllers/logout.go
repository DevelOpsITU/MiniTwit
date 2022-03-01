package controllers

import (
	"minitwit/functions"
	"minitwit/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
	functions.SetCookie(c, g)
	c.Redirect(http.StatusFound, "/")
}
