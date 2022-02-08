package main

import (
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
	g.user = false
	g.name = "jonas"

	out, err := tpl.Execute(gonja.Context{"first_name": "Christian", "last_name": "Mark", "g": g})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func main() {
	http.HandleFunc("/", examplePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
