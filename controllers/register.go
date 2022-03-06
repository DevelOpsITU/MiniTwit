package controllers

import (
	"minitwit/functions"
	"minitwit/log"
	"minitwit/logic"
	"minitwit/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
)

func registerHandlers(router *gin.Engine) {

	router.GET("/register", func(c *gin.Context) {
		handleRegisterGet(c.Writer, c)
	})
	router.POST("/register", func(c *gin.Context) {
		handleRegisterPost(c.Writer, c)
	})

}

var registerTemplate = gonja.Must(gonja.FromFile("templates/register.html"))

func handleRegisterGet(w gin.ResponseWriter, c *gin.Context) {
	user, err := functions.GetCookie(c)

	if err == nil && user.User.User_id != 0 {
		// If there are a cookie in the session i.e. no error when getting it
		c.Redirect(http.StatusFound, "/")
		return
	}

	out, err := registerTemplate.Execute(gonja.Context{"g": ""})
	if err != nil {
		log.Logger.Error().Err(err).Msg("Could not render register template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
	return
}

func handleRegisterPost(w gin.ResponseWriter, c *gin.Context) {
	cookieUser, err := functions.GetCookie(c)

	if err == nil && cookieUser.User.User_id != 0 {
		// If there are a cookie in the session i.e. no error when getting it
		c.Redirect(http.StatusFound, "/")
		return
	}

	registrationUser := models.RegistrationUser{
		Username:  c.PostForm("username"),
		Email:     c.PostForm("email"),
		Password1: c.PostForm("password"),
		Password2: c.PostForm("password2"),
	}

	err = logic.CreateUser(registrationUser)

	if err != nil {
		log.Logger.Error().Err(err).Str("username", registrationUser.Username).Str("email", registrationUser.Email).
			Msg("Could not register user")
		out, err := registerTemplate.Execute(gonja.Context{"g": "", "error": err.Error()})
		if err != nil {
			log.Logger.Error().Err(err).Msg("Could not render register template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(out))
		return
	}

	var g = models.Session{
		User:     models.User{}, //TODO: Maybe get the user from creation of the user
		Message:  true,
		Messages: []string{"You were successfully registered and can login now"},
	}
	functions.SetCookie(c, g)
	c.Redirect(http.StatusFound, "/login")
	return
}
