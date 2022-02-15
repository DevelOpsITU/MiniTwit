package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"minitwit/src/database"
	"minitwit/src/functions"
	"minitwit/src/models"
	"net/http"
	"strings"
)

func registerHandlers(router *gin.Engine) {

	router.GET("/register", func(c *gin.Context) {
		handleRegister(c.Writer, c.Request, c)
	})
	router.POST("/register", func(c *gin.Context) {
		handleRegister(c.Writer, c.Request, c)
	})

}

var registerTemplate = gonja.Must(gonja.FromFile("templates/register.html"))

func handleRegister(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	_, err := functions.GetCookie(c)

	if err == nil {
		// If there are a cookie in the session i.e. no error when getting it
		c.Redirect(http.StatusFound, "/")
		return
	}

	er := ""
	if r.Method == http.MethodPost {
		if c.PostForm("username") == "" {
			er = "You have to enter a username"
		} else if c.PostForm("email") == "" || !strings.Contains(c.PostForm("email"), "@") {
			er = "Your have to enter a valid email address"
		} else if c.PostForm("password") == "" {
			er = "You have to enter a password"
		} else if c.PostForm("password") != c.PostForm("password2") {
			er = "The two passwords do not match"
		} else if database.GetUserFromDb(c.PostForm("username")) != (models.User{}) {
			er = "The username is already taken"
		} else {
			database.AddUserToDb(c.PostForm("username"), c.PostForm("email"), c.PostForm("password"))
			var g = models.Session{
				User:     models.User{},
				Message:  true,
				Messages: []string{"You were successfully registered and can login now"},
			}
			data, _ := json.Marshal(g)
			c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}
	out, err := registerTemplate.Execute(gonja.Context{"g": "", "error": er})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}
