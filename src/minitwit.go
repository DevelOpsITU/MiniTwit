package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
)

func addNumbers(n1 int, n2 int) int {
	return n1 + n2
}

// Pre-compiling the templates at application startup using the
// little Must()-helper function (Must() will panic if FromFile()
// or FromString() will return with an error - that's it).
// It's faster to pre-compile it anywhere at startup and only
// execute the template later.

var tpl = gonja.Must(gonja.FromFile("templates/example.html"))

func examplePage(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request

	type User struct {
		username string
	}
	type structType struct {
		user User
	}
	var g structType
	//g.user.username = "jonas"

	out, err := tpl.Execute(gonja.Context{"first_name": "Christian", "last_name": "Mark", "g": g})
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
		examplePage(c.Writer, c.Request)
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
