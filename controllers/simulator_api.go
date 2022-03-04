package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"minitwit/logic"
	"minitwit/models"
	"net/http"
	"strconv"
)

var latest = -1

func simulationHandlers(router *gin.Engine) {

	// Register
	router.POST("/sim/register", func(c *gin.Context) {
		updateLatest(c)
		handleSimRegisterPost(c.Writer, c.Request)
	})

	// All Messages
	router.GET("/sim/msgs", func(c *gin.Context) {
		updateLatest(c)
		handleSimGetAllMessages(c)
	})

	// Latest
	router.GET("/sim/latest", func(c *gin.Context) {
		handleSimLatest(c.Writer)
	})

	// Message by user
	router.POST("/sim/msgs/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimAddMessage(c.Writer, c.Request, username)
	})

	router.GET("/sim/msgs/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimGetUserMessages(c.Writer, c, username)
	})

	// Follows
	router.GET("/sim/fllws/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimFollowUser(c.Writer, c.Request, c, username)
	})

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
func handleSimGetAllMessages(c *gin.Context) {
	/*err := not_req_from_simulator(r) // TODO: This should prob be added (he has it in his, but doesnt work when checking normally)
	if err != nil {
		fmt.Println(err)
		return
	}*/

	type MessageObj struct {
		Content string `json:"content"`
		PubDate int64  `json:"pub_date"`
		User    string `json:"user"`
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
		Latest int `json:"latest"`
	}
	LatestObj.Latest = latest

	js, _ := json.Marshal(LatestObj)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// DONE
func handleSimRegisterPost(w gin.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	Paylaod := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"pwd"`
	}{}
	err := decoder.Decode(&Paylaod)
	if err != nil {
		log.Println(err.Error())
		return
	}

	registrationUser := models.RegistrationUser{
		Username:  Paylaod.Username,
		Email:     Paylaod.Email,
		Password1: Paylaod.Password,
		Password2: Paylaod.Password,
	}

	err = logic.CreateUser(registrationUser)

	if err != nil {
		print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// Done
func handleSimAddMessage(w http.ResponseWriter, r *http.Request, username string) {

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	Payload := struct {
		Content string `json:"content"`
	}{}
	err := decoder.Decode(&Payload)
	if err != nil {
		log.Println(err)
		return
	}

	err = logic.AddMessageFromUsername(username, Payload.Content)

	if err != nil {
		log.Println("Could not add message: " + Payload.Content)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// Gets the user
func handleSimGetUserMessages(w http.ResponseWriter, c *gin.Context, username string) {

	var limit = 0
	var err error
	limitStr := c.Query("no")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Println("Had an input that was not a number i could parse. " + limitStr)
			limit = 9999999999
		}
	} else {
		limit = 9999999999 // Hacky but maybe it works for us.
	}

	twits, _, err := logic.GetUserTwits(username, limit)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type MessageObj struct {
		Content string `json:"content"`
		PubDate int64  `json:"pub_date"`
		User    string `json:"user"`
	}

	var msgsAsObject []MessageObj
	for _, msg := range twits {
		time, _ := strconv.Atoi(msg.Pub_date)
		msgObj := MessageObj{Content: msg.Text, PubDate: int64(time), User: msg.Username}
		msgsAsObject = append(msgsAsObject, msgObj)
	}
	js, _ := json.Marshal(msgsAsObject)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(js)
}

// Done
func handleSimFollowUser(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {

	var limit = 0
	var err error
	limitStr := c.Query("no")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Println("Had an input that was not a number i could parse. " + limitStr)
			limit = 9999999999
		}
	} else {
		limit = 9999999999 // Hacky but maybe it works for us.
	}

	if r.Method == "POST" {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		Payload := struct {
			Follow   string `json:"follow"`
			Unfollow string `json:"unfollow"`
		}{}
		err := decoder.Decode(&Payload)
		if err != nil {
			log.Println(err)
			return
		}

		if Payload.Follow != "" {
			followUsername := Payload.Follow

			err = logic.FollowUserFromUsername(username, followUsername)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else if Payload.Unfollow != "" {
			unfollowUsername := Payload.Unfollow

			err = logic.UnFollowUserFromUsername(username, unfollowUsername)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	} else if r.Method == "GET" {

		followedByUser, err := logic.GetUserFollowerUsernames(username, limit)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		type followsObj struct {
			Follows []string `json:"follows"`
		}
		var follows followsObj

		for _, username := range followedByUser {
			follows.Follows = append(follows.Follows, username)
		}

		if followedByUser == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		usersAsJson, _ := json.Marshal(follows)

		w.Header().Set("Content-Type", "application/json")
		w.Write(usersAsJson)
		w.WriteHeader(http.StatusOK)
	}
	return

}
