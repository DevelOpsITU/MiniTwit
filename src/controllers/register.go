package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"
)

func registerHandlers(router *gin.Engine) {

	router.GET("/register", func(c *gin.Context) {
		handleRegisterGet(c.Writer, c.Request, c)
	})
	router.POST("/register", func(c *gin.Context) {
		handleRegisterPost(c.Writer, c.Request, c)
	})

}

var registerTemplate = gonja.Must(gonja.FromFile("templates/register.html"))

func handleRegisterGet(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	user, err := functions.GetCookie(c)

	if err == nil && user.User.User_id != 0 {
		// If there are a cookie in the session i.e. no error when getting it
		c.Redirect(http.StatusFound, "/")
		return
	}

	out, err := registerTemplate.Execute(gonja.Context{"g": ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func handleRegisterPost(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
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
