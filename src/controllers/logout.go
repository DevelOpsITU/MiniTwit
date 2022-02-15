package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minitwit/src/models"
	"net/http"
)

func HandleLogout(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
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
