package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"minitwit/src/database"
	"minitwit/src/models"
	"time"

	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/noirbizarre/gonja"
	"golang.org/x/crypto/pbkdf2"
)

/****************************************
*			REST OF PROGRAM				*
****************************************/

func getCookie(c *gin.Context) (models.Session, error) {
	var g models.Session
	cookie, err := c.Cookie("session")

	// If there is no cookie
	if err != nil {
		return g, err
	} else {
		//data,_ := json.Marshal(g)
		//c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
		json.Unmarshal([]byte(cookie), &g)
		//println("Found Cookie:", string([]byte(cookie)))

	}
	newCookie := g
	newCookie.Message = false
	newCookie.Messages = nil
	setCookie(c, newCookie)

	return g, nil

}

func setCookie(c *gin.Context, session models.Session) {

	data, _ := json.Marshal(session)
	c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
}

// Pre-compiling the templates at application startup using the
// little Must()-helper function (Must() will panic if FromFile()
// or FromString() will return with an error - that's it).
// It's faster to pre-compile it anywhere at startup and only
// execute the template later.

var timelineTemplate = gonja.Must(gonja.FromFile("templates/timeline.html"))
var loginTemplate = gonja.Must(gonja.FromFile("templates/login.html"))
var registerTemplate = gonja.Must(gonja.FromFile("templates/register.html"))

var g models.Session

// Route /
func handleTimeline(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// Execute the template per HTTP request
	request := getEndpoint(r)
	data, err := getCookie(c)
	g = data

	// If there is no cookie
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	//set g = "None" if g.user should return false in jinja

	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func handleUnFollowUser(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {
	data, err := getCookie(c)
	g = data

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	whom_id := database.GetUserFromDb(username)

	// TODO: check if followed before trying this
	db := database.ConnectDb()

	query, err := db.Prepare("DELETE FROM follower WHERE who_id = ? AND whom_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(g.User.User_id, whom_id.User_id)

	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	// set Message in cookie
	cookie := models.Session{
		User:     g.User,
		Message:  true,
		Messages: []string{"You are no longer following " + username},
	}
	newdata, _ := json.Marshal(cookie)
	c.SetCookie("session", string(newdata), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}

func handleFollowUser(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {
	data, err := getCookie(c)
	g = data

	// If there is no cookie / no user logged in
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	whom_id := database.GetUserFromDb(username)

	db := database.ConnectDb()

	query, err := db.Prepare("INSERT INTO follower (who_id, whom_id) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(g.User.User_id, whom_id.User_id)

	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	// set Message in cookie
	cookie := models.Session{
		User:     g.User,
		Message:  true,
		Messages: []string{"You are now following " + username},
	}
	newdata, _ := json.Marshal(cookie)
	c.SetCookie("session", string(newdata), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}

func handleAddMessage(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	data, err := getCookie(c)
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

func handleUserTimeline(w http.ResponseWriter, r *http.Request, c *gin.Context, username string) {

	user := database.GetUserFromDb(username)

	messages := database.GetUserMessages(user.User_id)
	twits := CovertMessagesToTwits(&messages)
	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": r, "messages": twits})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func getEndpoint(r *http.Request) models.Request {
	var request = models.Request{r.URL.Path}

	if request.Endpoint == "/public" {
		request.Endpoint = "public_timeline"
	} else if len(request.Endpoint) > 1 {
		request.Endpoint = "user_timeline"
	} else {
		request.Endpoint = ""
	}
	return request
}

func handlePublicTimeline(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	request := getEndpoint(r)

	messages := database.GetAllMessages()

	twits := CovertMessagesToTwits(&messages)
	//print(string(request))
	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request, "messages": twits})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func CovertMessagesToTwits(messages *[]models.Message) []models.Twit {
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

func handleLogin(w gin.ResponseWriter, r *http.Request, c *gin.Context) {

	//var error = ""
	//GetAllMessages()

	if r.Method == http.MethodPost {
		username := c.PostForm("username")
		pw := c.PostForm("password")
		print(pw)
		print(username)
		user := database.GetUserFromDb(username)
		print(user.Pw_hash)
		s := strings.Split(user.Pw_hash, ":")

		s2 := strings.Split(s[2], "$")

		pwIteration := s2[0]
		pwSalt := s2[1]
		pwHash := s2[2]
		fmt.Println(pwIteration)
		fmt.Println(pwSalt)
		fmt.Println(pwHash)

		//salt := []byte(user.pw_hash[21:37]) // TODO: what is this used for?
		pwIteration_int, _ := strconv.Atoi(pwIteration)

		dk := pbkdf2.Key([]byte(pw), []byte(pwSalt), pwIteration_int, 32, sha256.New)

		fmt.Printf("\nsha256: %x\n", []byte(dk))

		// TODO: same as l.293 (What is thi used for?)
		//fmt.Printf("salt: %x\n", string(salt))
		//fmt.Println("len(salt)", len(salt),
		//	"\nlen(hashed)", len(dk))

		if hex.EncodeToString(dk) != pwHash {
			// Invalid
			print("Invalid password")
			c.Redirect(http.StatusFound, "/login")

			return
		} else {
			// User is authenticated

			var g = models.Session{
				User:     user,
				Message:  true,
				Messages: []string{"You were successfully logged in"},
			}
			data, _ := json.Marshal(g)
			c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusFound, "/")
			return
		}

	}
	cookie, err := getCookie(c)
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

func handleRegister(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	_, err := getCookie(c)

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

func handleLogout(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	// reset cookie
	g := models.Session{
		User:     models.User{},
		Message:  true,
		Messages: []string{"You were logged out"},
	}
	data, _ := json.Marshal(g)
	c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/")
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pongA",
		})
	})

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		handleTimeline(c.Writer, c.Request, c)
	})
	router.GET("/public", func(c *gin.Context) {
		handlePublicTimeline(c.Writer, c.Request, c)
	})

	router.GET("/login", func(c *gin.Context) {
		handleLogin(c.Writer, c.Request, c)
	})

	router.POST("/login", func(c *gin.Context) {
		handleLogin(c.Writer, c.Request, c)
	})

	// Logout
	router.GET("/logout", func(c *gin.Context) {
		handleLogout(c.Writer, c.Request, c)
	})

	// Register
	router.GET("/register", func(c *gin.Context) {
		handleRegister(c.Writer, c.Request, c)
	})
	router.POST("/register", func(c *gin.Context) {
		handleRegister(c.Writer, c.Request, c)
	})

	// User timeline
	router.GET("/:user", func(c *gin.Context) {
		username := c.Param("user")
		handleUserTimeline(c.Writer, c.Request, c, username)
	})

	// Follow
	router.GET("/:user/follow", func(c *gin.Context) {
		username := c.Param("user")
		handleFollowUser(c.Writer, c.Request, c, username)
	})

	// Unfollow
	router.GET("/:user/unfollow", func(c *gin.Context) {
		username := c.Param("user")
		handleUnFollowUser(c.Writer, c.Request, c, username)
	})

	// Add message
	router.GET("/add_message", func(c *gin.Context) {
		handleAddMessage(c.Writer, c.Request, c)
	})

	router.LoadHTMLFiles("./src/test.html")

	/*
	 FOR TESTING GO TOOL 'FRESH': 'go install github.com/pilu/fresh'
	 TRY TO RUN COMMAND: 'fresh -c my_fresh_runner.conf' AND
	 THEN MAKE CHANGES TO THE 'test.html' OR 'minitwit.go' FILES.
	 IF NO ERROR, THEN FRESH SHOULD BUILD AND RUN THE 'minitwit.go' CODE.
	 THE CHANGES SHOULD BE SEEN REFLECTED ON 'http://localhost:8080/test/test.html'.

	 OBS: MAYBE TURN OFF AUTO-SAVING, SO STUFF IS ONLY BUILD AND RAN, WHEN YOU WANT IT TO.
	*/
	router.Static("/test", "./src")

	router.Run(":8080")
	//router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
