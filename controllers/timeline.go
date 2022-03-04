package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"minitwit/functions"
	"minitwit/log"
	"minitwit/logic"
	"minitwit/models"
	"net/http"
)

func timelineHandlers(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		handleRootTimeline(c.Writer, c.Request, c)
	})
	router.GET("/public", func(c *gin.Context) {
		handlePublicTimeline(c.Writer, c.Request)
	})

}

var timelineTemplate = gonja.Must(gonja.FromFile("templates/timeline.html"))

var g models.Session

func handleUserTimeline(w http.ResponseWriter, r *http.Request, username string) {

	request := functions.GetEndpoint(r)

	twits, user, err := logic.GetUserTwits(username, 30)
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

func handleRootTimeline(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// Execute the template per HTTP request
	request := functions.GetEndpoint(r)
	data, err := functions.GetCookie(c)
	g = data

	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	twits, err := logic.GetPersonalTimelineTwits(g.User)

	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits, "profile_user": g.User})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func handlePublicTimeline(w gin.ResponseWriter, r *http.Request) {
	request := functions.GetEndpoint(r)

	twits, err := logic.GetPublicTimelineTwits()

	if err != nil {
		log.Logger.Error().Err(err).Caller().Msg("Could not get public messages")

	}
	//print(string(request))
	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits})
	if err != nil {
		log.Logger.Error().Err(err).Caller().Msg("Could not generate the timeline template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}
