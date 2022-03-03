package controllers

import (
	"encoding/json"
	"minitwit/functions"
	"minitwit/logic"
	"minitwit/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func userHandlers(router *gin.Engine) {

	router.GET("/:user/follow", func(c *gin.Context) {
		username := c.Param("user")
		handleFollowUser(c.Writer, c, username)
	})

	router.GET("/:user/unfollow", func(c *gin.Context) {
		username := c.Param("user")
		handleUnFollowUser(c.Writer, c, username)
	})

	// User timeline
	router.GET("/:user", func(c *gin.Context) {
		username := c.Param("user")
		handleUserTimeline(c.Writer, c.Request, username)
	})

}

func handleUnFollowUser(w http.ResponseWriter, c *gin.Context, username string) {
	data, err := functions.GetCookie(c)
	g = data

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		response := make(map[string]string)
		response["message"] = "401 Unautherized - no user logged in!"
		jsonText, err := json.Marshal(response)
		if err != nil {
			panic("Error handling JSON marshal")
		}
		w.Write(jsonText)
		return
		// c.Redirect(http.StatusFound, "/public")
	}

	logic.UnFollowUser(g.User.User_id, username)

	// set Message in cookie
	g := models.Session{
		User:     g.User,
		Message:  true,
		Messages: []string{"You are no longer following " + username},
	}
	functions.SetCookie(c, g)
	c.Redirect(http.StatusFound, "/")
}

func handleFollowUser(w http.ResponseWriter, c *gin.Context, username string) {
	data, err := functions.GetCookie(c)
	g = data

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		response := make(map[string]string)
		response["message"] = "401 Unautherized - no user logged in!"
		jsonText, err := json.Marshal(response)
		if err != nil {
			panic("Error handling JSON marshal")
		}
		w.Write(jsonText)
		return
		// c.Redirect(http.StatusFound, "/public")
	}

	logic.FollowUser(g.User.User_id, username)

	// set Message in cookie
	g := models.Session{
		User:     g.User,
		Message:  true,
		Messages: []string{"You are now following " + username},
	}
	functions.SetCookie(c, g)
	c.Redirect(http.StatusFound, "/")
}
