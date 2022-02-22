package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"minitwit/src/database"
	"minitwit/src/functions"
	"minitwit/src/logic"
	"minitwit/src/models"
	"net/http"
	"strconv"
)

var latest = -1

func simulationHandlers(router *gin.Engine) {

	// Register // TODO: NOT DONE
	router.POST("/sim/register", func(c *gin.Context) {
		updateLatest(c)
		handleSimRegisterPost(c.Writer, c.Request, c)
	})

	// All Messages // TODO: NOT DONE
	router.GET("/sim/msgs", func(c *gin.Context) {
		updateLatest(c)
		handleSimGetAllMessages(c, c.Request)
	})

	// Latest
	router.GET("/sim/latest", func(c *gin.Context) {
		handleSimLatest(c.Writer)
	})

	// Message by user
	router.POST("/sim/msgs/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimAddMessage(c.Writer, c.Request, c, username)
	})
	router.GET("/sim/msgs/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimGetUserMessages(c.Writer, c.Request, c, username)
	})

	// Follows
	router.GET("/sim/fllws/:username", func(c *gin.Context) {
		updateLatest(c)
		username := c.Param("username")
		handleSimFollowUser(c.Writer, c.Request, c, username)
	})

}

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

// TODO: Not done
func handleSimRegisterPost(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	cookieUser, err := functions.GetCookie(c)

	if err == nil && cookieUser.User.User_id != 0 {
		// If there are a cookie in the session i.e. no error when getting it
		c.Redirect(http.StatusNoContent, "/")
		return
	}

	registationUser := models.RegistrationUser{
		c.PostForm("username"),
		c.PostForm("email"),
		c.PostForm("password"),
		c.PostForm("password2"),
	}

	err = logic.CreateUser(registationUser)

	if err != nil {
		out, err := registerTemplate.Execute(gonja.Context{"g": "", "error": err.Error()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(out))
		return
	} else {
		c.Redirect(http.StatusNoContent, "/")
		return
	}
}

// TODO: Not done
func handleSimAddMessage(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {

	g, err := functions.GetCookie(c)

	/*x := struct {
		content 	string
	}*/
	content := r.Form.Get("content")
	user, err := database.GetUserFromDb(username)

	err = logic.AddMessage(user, content)

	if err != nil {
		println(err.Error())
	} else {
		g = models.Session{
			User:     g.User,
			Message:  true,
			Messages: []string{"Your message was recorded"},
		}
	}

	//var data, _ = json.Marshal(g)
	c.Redirect(http.StatusOK, "/")
}

// TODO: Not done
func handleSimGetUserMessages(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {

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

// TODO: not done
func handleSimFollowUser(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {

	user, _ := database.GetUserFromDb(username)

	if r.Method == "POST" {
		if r.Form.Get("follow") != "" {
			followUsername := r.Form.Get("follow")

			logic.FollowUser(user.User_id, followUsername)
		} else if r.Form.Get("unfollow") != "" {
			unfollowUsername := r.Form.Get("unfollow")

			logic.UnFollowUser(user.User_id, unfollowUsername)
		}

		c.Redirect(http.StatusNoContent, "/")
	} else if r.Method == "GET" {
		followedByUser := logic.GetUsernameOfWhoFollowsUser(user.User_id, c.Query("no"))
		usersAsJson, _ := json.Marshal(followedByUser)

		w.Write(usersAsJson)
	}

}
