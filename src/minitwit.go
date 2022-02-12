package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"net/http"
)

type User struct {
	Username string
}
type Session struct {
	User     User
	Message  bool
	Messages []string
}

type Request struct {
	Endpoint string
}

func getCookie(c *gin.Context) (Session, error) {
	var g Session
	cookie, err := c.Cookie("session")

	// If there is no cookie
	if err != nil {
		return g, err
	} else {
		//data,_ := json.Marshal(g)
		//c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
		json.Unmarshal([]byte(cookie), &g)
		print("Found Cookie:", string([]byte(cookie)))
		return g, nil
	}
}

// Pre-compiling the templates at application startup using the
// little Must()-helper function (Must() will panic if FromFile()
// or FromString() will return with an error - that's it).
// It's faster to pre-compile it anywhere at startup and only
// execute the template later.

var tpl = gonja.Must(gonja.FromFile("templates/timeline.html"))

// Route /
func handleTimeline(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// Execute the template per HTTP request

	var g Session
	g, err := getCookie(c)

	// If there is no cookie
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	//set g = "None" if g.user should return false in jinja

	out, err := tpl.Execute(gonja.Context{"first_name": "Christian", "last_name": "Mark", "g": g})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func handlePublicTimeline(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	var request = Request{r.URL.Path}

	if request.Endpoint == "/public" {
		request.Endpoint = "public_timeline"
	} else if len(request.Endpoint) > 1 {
		request.Endpoint = "user_timeline"
	} else {
		request.Endpoint = ""
	}

	type Message struct {
		Email    string
		Username string
		Pub_date string
		Text     string
	}

	var messages = []Message{{"User@email.com", "user1", "dato", "Twit1"}, {"User2@email.com", "user2", "dato", "Twit2"}}
	//print(string(request))
	out, err := tpl.Execute(gonja.Context{"first_name": "Christian", "last_name": "Mark", "g": "", "request": request, "messages": messages})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
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
