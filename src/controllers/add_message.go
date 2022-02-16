package controllers

import (
	"encoding/json"
	"log"
	"minitwit/src/database"
	"minitwit/src/functions"
	"minitwit/src/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func addMessageHandlers(router *gin.Engine) {

	// Add message
	router.POST("/add_message", func(c *gin.Context) {
		handleAddMessage(c.Writer, c.Request, c)
	})

}

func handleAddMessage(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	data, err := functions.GetCookie(c)
	g = data

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	// Insert message into db
	if c.PostForm("text") != "" {
		db := database.ConnectDb()

		query, err := db.Prepare(`INSERT INTO message (author_id, text, pub_date, flagged) 
			VALUES (?, ?, ?, 0)`)

		if err != nil {
			log.Fatal(err)
		}
		_, err = query.Exec(g.User.User_id, c.PostForm("text"), time.Now().Unix())

		if err != nil {
			log.Fatal(err)
		}
		defer query.Close()

		var g = models.Session{
			User:     g.User,
			Message:  true,
			Messages: []string{"Your message was recorded"},
		}
		data, _ := json.Marshal(g)
		c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/")
	}
	c.Redirect(http.StatusFound, "/")
}
