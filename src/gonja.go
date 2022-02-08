package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"net/http"
)

// Pre-compiling the templates at application startup using the
// little Must()-helper function (Must() will panic if FromFile()
// or FromString() will return with an error - that's it).
// It's faster to pre-compile it anywhere at startup and only
// execute the template later.
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

	out, err := tpl.Execute(gonja.Context{"first_name": "Christian", "last_name": "Mark", "g": g})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		examplePage(c.Writer, c.Request)
	})

	router.Run(":8080")
	//router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
