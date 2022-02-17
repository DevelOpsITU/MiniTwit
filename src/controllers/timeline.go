package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"minitwit/src/database"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"
)

func timelineHandlers(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		handleTimeline(c.Writer, c.Request, c)
	})
	router.GET("/public", func(c *gin.Context) {
		handlePublicTimeline(c.Writer, c.Request, c)
	})

}

var timelineTemplate = gonja.Must(gonja.FromFile("templates/timeline.html"))

var g models.Session

func handleUserTimeline(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {

	request := functions.GetEndpoint(r)

	twits, user, err := logic.GetUserTwits(username)
	if err != nil {
		http.Error(w, "404 - "+err.Error(), http.StatusNotFound)
		return
	}

	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits, "profile_user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(out))
}

func handleTimeline(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// Execute the template per HTTP request
	request := functions.GetEndpoint(r)
	data, err := functions.GetCookie(c)
	g = data

	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	user, _ := database.GetUserFromDb(g.User.Username)
	messages := database.GetUserMessages(g.User.User_id)
	twits := logic.ConvertMessagesToTwits(&messages)

	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits, "profile_user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func handlePublicTimeline(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	request := functions.GetEndpoint(r)

	messages := database.GetAllMessages()

	twits := logic.ConvertMessagesToTwits(&messages)
	//print(string(request))
	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}
