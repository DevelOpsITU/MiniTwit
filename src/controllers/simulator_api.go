package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"
	"strconv"
)

var latest = -1

func simulationHandlers(router *gin.Engine) {

	// Register
	router.POST("/sim/register", func(c *gin.Context) {
		latest, _ = strconv.Atoi(c.Request.URL.Query().Get("latest"))
		handleSimRegisterPost(c.Writer, c.Request, c)
	})

	// All Messages
	router.GET("/sim/msgs", func(c *gin.Context) {
		latest, _ = strconv.Atoi(c.Request.URL.Query().Get("latest"))
		handleAllMessages(c.Writer)
	})

	// Latest
	router.GET("/sim/latest", func(c *gin.Context) {
		handleLatest(c.Writer)
	})

}

func handleAllMessages(w gin.ResponseWriter) {
	twits, _ := logic.GetPublicTimelineTwits()
	js, _ := json.Marshal(twits)

	w.Write(js)
}

func handleLatest(w gin.ResponseWriter) {
	json.Marshal(latest)
	js, _ := json.Marshal(latest)
	w.Write(js)
}

func handleSimRegisterPost(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	cookieUser, err := functions.GetCookie(c)

	if err == nil && cookieUser.User.User_id != 0 {
		// If there are a cookie in the session i.e. no error when getting it
		c.Redirect(http.StatusFound, "/")
		return
	}

	registationUser := models.RegistrationUser{
		c.PostForm("username"),
		c.PostForm("email"),
		c.PostForm("password"),
		c.PostForm("password2"),
	}

	err = logic.CreateUser(registationUser)

	if err != nil {
		out, err := registerTemplate.Execute(gonja.Context{"g": "", "error": err.Error()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(out))
		return
	} else {
		var g = models.Session{
			User:     models.User{}, //TODO: Maybe get the user from creation of the user
			Message:  true,
			Messages: []string{"You were successfully registered and can login now"},
		}
		data, _ := json.Marshal(g)
		c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/login")
		return
	}
}
