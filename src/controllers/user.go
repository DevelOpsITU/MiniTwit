package controllers

import (
	"encoding/json"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func userHandlers(router *gin.Engine) {

	router.GET("/:user/follow", func(c *gin.Context) {
		username := c.Param("user")
		handleFollowUser(c.Writer, c.Request, c, username)
	})

	router.GET("/:user/unfollow", func(c *gin.Context) {
		username := c.Param("user")
		handleUnFollowUser(c.Writer, c.Request, c, username)
	})

	// User timeline
	router.GET("/:user", func(c *gin.Context) {
		username := c.Param("user")
		handleUserTimeline(c.Writer, c.Request, c, username)
	})

}

func handleUnFollowUser(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {
	data, err := functions.GetCookie(c)
	g = data

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	logic.UnFollowUser(g.User.User_id, username)

	// set Message in cookie
	cookie := models.Session{
		User:     g.User,
		Message:  true,
		Messages: []string{"You are no longer following " + username},
	}
	newdata, _ := json.Marshal(cookie)
	c.SetCookie("session", string(newdata), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}

func handleFollowUser(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {
	data, err := functions.GetCookie(c)
	g = data

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	logic.FollowUser(g.User.User_id, username)

	// set Message in cookie
	cookie := models.Session{
		User:     g.User,
		Message:  true,
		Messages: []string{"You are now following " + username},
	}
	newdata, _ := json.Marshal(cookie)
	c.SetCookie("session", string(newdata), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}
