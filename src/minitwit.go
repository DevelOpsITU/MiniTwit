package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"net/http"
)

func addNumbers(n1 int, n2 int) int {
	return n1 + n2
}

// Pre-compiling the templates at application startup using the
// little Must()-helper function (Must() will panic if FromFile()
// or FromString() will return with an error - that's it).
// It's faster to pre-compile it anywhere at startup and only
// execute the template later.

var tpl_append = gonja.Must(gonja.FromFile("templates/append.html"))
var tpl = gonja.Must(gonja.FromFile("templates/example.html"))

func examplePage(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	type structType struct {
		user bool
		name string
	}
	var g structType
	g.user = true
	g.name = "jonas"

	out2, err2 := tpl.Execute(gonja.Context{"first_name": "Christian", "last_name": "Mark", "g": g})
	out, err := tpl_append.Execute(gonja.Context{"first_name": "Christian", "last_name": "Mark", "g": g})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	print(out, err2)
	w.Write([]byte(out2))
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		examplePage(c.Writer, c.Request)
	})

	router.Run(":8080")
	//router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
