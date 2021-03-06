package controllers

import (
	"minitwit/functions"
	"minitwit/log"
	"minitwit/logic"
	"minitwit/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addMessageHandlers(router *gin.Engine) {

	// Add message
	router.POST("/add_message", func(c *gin.Context) {
		handleAddMessage(c)
	})

}

func handleAddMessage(c *gin.Context) {

	g, err := functions.GetCookie(c)

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	err = logic.AddMessageFromUserModel(g.User, c.PostForm("text"))

	if err != nil {
		log.Logger.Error().Err(err).Str("text", c.PostForm("text")).Msg("Could not add message")
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
