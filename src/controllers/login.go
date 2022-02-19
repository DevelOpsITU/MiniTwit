package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"
)

func loginHandlers(router *gin.Engine) {

	router.GET("/login", func(c *gin.Context) {
		handleloginGet(c.Writer, c.Request, c)
	})

	router.POST("/login", func(c *gin.Context) {
		handleLogin(c.Writer, c.Request, c)
	})
}

var loginTemplate = gonja.Must(gonja.FromFile("templates/login.html"))

func handleloginGet(w gin.ResponseWriter, r *http.Request, c *gin.Context) {

	cookie, err := functions.GetCookie(c)
	if err != nil {
		out, err := loginTemplate.Execute(gonja.Context{"g": ""})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(out))
	} else {
		out, err := loginTemplate.Execute(gonja.Context{"g": cookie})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(out))
	}

}

func handleLogin(w gin.ResponseWriter, r *http.Request, c *gin.Context) {

	username := c.PostForm("username")
	pw := c.PostForm("password")

	user, err := logic.CheckPassword(username, pw)

	if err != nil {
		println(err.Error())
		var g = models.Session{
			User:     models.User{},
			Message:  true,
			Messages: []string{err.Error()}}
		functions.SetCookie(c, g)
		c.Redirect(http.StatusFound, "/login")
		return
	} else {
		var g = models.Session{
			User:     user,
			Message:  true,
			Messages: []string{"You were successfully logged in"},
		}
		functions.SetCookie(c, g)
		c.Redirect(http.StatusFound, "/")
		return

	}
}