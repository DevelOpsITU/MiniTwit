package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"minitwit/src/database"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
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

	messages, user, err := logic.GetUserMessages(username)
	if err != nil {
		http.Error(w, "404 - "+err.Error(), http.StatusNotFound)
		return
	}

	twits := ConvertMessagesToTwits(&messages)
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
	twits := ConvertMessagesToTwits(&messages)

	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits, "profile_user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func handlePublicTimeline(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	request := functions.GetEndpoint(r)

	messages := database.GetAllMessages()

	twits := ConvertMessagesToTwits(&messages)
	//print(string(request))
	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func ConvertMessagesToTwits(messages *[]models.Message) []models.Twit {
	var twits []models.Twit
	for _, message := range *messages {
		twits = append(twits, models.Twit{getGavaterUrl(message.Email, 48), message.Username, strconv.Itoa(int(message.Pubdate)), message.Text})
	}
	print(twits)
	return twits

}

func getGavaterUrl(email string, size int) string {
	data := []byte(strings.ToLower(strings.TrimSpace(email)))
	hash := md5.Sum(data)
	hashStr := hex.EncodeToString(hash[:])

	str := []string{"http://www.gravatar.com/avatar/", hashStr, "?d=identicon&s=", strconv.Itoa(size)}
	return strings.Join(str, "")
}
