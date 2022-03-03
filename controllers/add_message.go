package controllers

import (
	"minitwit/functions"
	"minitwit/logic"
	"minitwit/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addMessageHandlers(router *gin.Engine) {

	// Add message
	router.POST("/add_message", func(c *gin.Context) {
		handleAddMessage(c.Writer, c.Request, c)
	})

}

func handleAddMessage(w http.ResponseWriter, r *http.Request, c *gin.Context) {

	g, err := functions.GetCookie(c)

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	err = logic.AddMessage(g.User, c.PostForm("text"))

	if err != nil {
		println(err.Error())
	} else {
		g = models.Session{
			User:     g.User,
			Message:  true,
			Messages: []string{"Your message was recorded"},
		}
	}

	functions.SetCookie(c, g)
	c.Redirect(http.StatusFound, "/")
}
