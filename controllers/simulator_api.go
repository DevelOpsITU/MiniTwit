package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"minitwit/database"
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

	//content := r.Form.Get("content")
	/*tmp, _ := c.GetRawData()
	type messagDTO = struct {
		content string
	}
	var payload messagDTO
	json.Unmarshal(tmp, &payload)
	*/

	print(Payload.Content)
	user, _ := database.GetUserFromDb(username)
	_ = logic.AddMessage(user, Payload.Content)

	w.WriteHeader(http.StatusNoContent)
}

// Done
func handleSimGetUserMessages(w http.ResponseWriter, c *gin.Context, username string) {
	user, err := database.GetUserFromDb(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages := logic.GetUserSimulationMessages(user, c.Query("no"))

	type MessageObj struct {
		Content string `json:"content"`
		PubDate int64  `json:"pub_date"`
		User    string `json:"user"`
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
		return
	}

	if r.Method == "POST" {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		Paylaod := struct {
			Follow   string `json:"follow"`
			Unfollow string `json:"unfollow"`
		}{}
		err = decoder.Decode(&Paylaod)
		if err != nil {
			log.Println(err)
			return
		}

		if Paylaod.Follow != "" {
			followUsername := Paylaod.Follow
			followUser, err := database.GetUserFromDb(followUsername)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			err = logic.FollowSimulationUser(user.User_id, followUser)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else if Paylaod.Unfollow != "" {
			unfollowUsername := Paylaod.Unfollow

			unfollowUser, err := database.GetUserFromDb(unfollowUsername)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)

			}

			logic.UnFollowSimulationUser(user.User_id, unfollowUser)
		}

		w.WriteHeader(http.StatusNoContent)
	} else if r.Method == "GET" {
		followedByUser := logic.GetUsernameOfWhoFollowsUser(user.User_id, c.Query("no"))

		type followsObj struct {
			Follows []string `json:"follows"`
		}
		var follows followsObj

		for user_ := range followedByUser {
			follows.Follows = append(follows.Follows, followedByUser[user_])
		}

		if followedByUser == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		usersAsJson, _ := json.Marshal(follows)

		w.Header().Set("Content-Type", "application/json")
		w.Write(usersAsJson)
	}
	return

}
