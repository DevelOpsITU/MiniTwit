package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"minitwit/src/database"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"
	"strconv"
)

var latest = -1

func simulationHandlers(router *gin.Engine) {

	// Register // TODO: DONE
	router.POST("/sim/register", func(c *gin.Context) {
		updateLatest(c)
		handleSimRegisterPost(c.Writer, c.Request, c)
	})

	// All Messages // TODO: DONE
	router.GET("/sim/msgs", func(c *gin.Context) {
		updateLatest(c)
		handleSimGetAllMessages(c, c.Request)
	})

	// Latest // TODO: DONE
	router.GET("/sim/latest", func(c *gin.Context) {
		handleSimLatest(c.Writer)
	})

	// Message by user // TODO: DONE
	router.POST("/sim/msgs/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimAddMessage(c.Writer, c.Request, c, username)
	})
	// TODO: DONE
	router.GET("/sim/msgs/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimGetUserMessages(c.Writer, c.Request, c, username)
	})

	// Follows // TODO: DONE
	router.GET("/sim/fllws/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimFollowUser(c.Writer, c.Request, c, username)
	})
	// TODO: DONE
	router.POST("/sim/fllws/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimFollowUser(c.Writer, c.Request, c, username)
	})

}

// TODO: MAYBE USE THIS
func not_req_from_simulator(r *http.Request) error {
	fromSimulator := r.Header.Get("Authorization")
	if fromSimulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		return errors.New("You are not authorized to use this resource!")
	}
	return nil
}

func updateLatest(c *gin.Context) {
	atoi, err := strconv.Atoi(c.Query("latest"))
	if err == nil {
		latest = atoi
	}
}

// Done
func handleSimGetAllMessages(c *gin.Context, r *http.Request) {
	/*err := not_req_from_simulator(r) // TODO: This should prob be added (he has it in his, but doesnt work when checking normally)
	if err != nil {
		fmt.Println(err)
		return
	}*/

	type MessageObj struct {
		Content string `json:"content`
		PubDate int64  `json:"pub_date`
		User    string `json:"user`
	}
	messages := logic.GetAllSimulationMessages(c.Query("no"))

	var msgsAsObject []MessageObj
	for _, msg := range messages {
		msgObj := MessageObj{Content: msg.Text, PubDate: msg.Pubdate, User: msg.Username}
		msgsAsObject = append(msgsAsObject, msgObj)
	}
	js, _ := json.Marshal(msgsAsObject)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(js)
}

// DONE
func handleSimLatest(w gin.ResponseWriter) {

	var LatestObj struct {
		Latest int `json:"latest`
	}
	LatestObj.Latest = latest

	js, _ := json.Marshal(LatestObj)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// DONE
func handleSimRegisterPost(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	cookieUser, err := functions.GetCookie(c)

	if err == nil && cookieUser.User.User_id != 0 {
		// If there are a cookie in the session i.e. no error when getting it
		w.WriteHeader(http.StatusNoContent)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// Done
func handleSimAddMessage(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {
	content := r.Form.Get("content")
	user, _ := database.GetUserFromDb(username)
	_ = logic.AddMessage(user, content)

	w.WriteHeader(http.StatusNoContent)
}

// Done
func handleSimGetUserMessages(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {
	user, err := database.GetUserFromDb(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages := logic.GetUserSimulationMessages(user, c.Query("no"))

	type MessageObj struct {
		Content string `json:"content`
		PubDate int64  `json:"pub_date`
		User    string `json:"user`
	}

	var msgsAsObject []MessageObj
	for _, msg := range messages {
		msgObj := MessageObj{Content: msg.Text, PubDate: msg.Pubdate, User: msg.Username}
		msgsAsObject = append(msgsAsObject, msgObj)
	}
	js, _ := json.Marshal(msgsAsObject)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(js)
}

// Done
func handleSimFollowUser(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {

	user, err := database.GetUserFromDb(username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if r.Method == "POST" {
		if r.Form.Get("follow") != "" {
			followUsername := r.Form.Get("follow")
			followUser, err := database.GetUserFromDb(followUsername)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			logic.FollowSimulationUser(user.User_id, followUser)
		} else if r.Form.Get("unfollow") != "" {
			unfollowUsername := r.Form.Get("unfollow")

			unfollowUser, err := database.GetUserFromDb(unfollowUsername)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			logic.UnFollowSimulationUser(user.User_id, unfollowUser)
		}

		w.WriteHeader(http.StatusNoContent)
		return
	} else if r.Method == "GET" {
		followedByUser := logic.GetUsernameOfWhoFollowsUser(user.User_id, c.Query("no"))
		usersAsJson, _ := json.Marshal(followedByUser)

		w.Header().Set("Content-Type", "application/json")
		w.Write(usersAsJson)
	}

}
