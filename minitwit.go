package main

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"net/http"

	// "github.com/flosch/pongo2/v5" // antivirus does not seem to like this (!OBS known problem on windows for go https://go.dev/doc/faq#virus)
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//const DATABASE = "tmp/minitwit.db"
const DATABASE = "C:/Users/hardk/source/repos/MiniTwit/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

/****************************************
*		   DATABASE ENTITIES			*
****************************************/
type Message struct {
	MessageId int
	AuthorId  int
	Text      string
	Pubdate   int64
	Flagged   bool
}

/****************************************
*			DATABASE RELATED			*
****************************************/
func connectDb() *sql.DB {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}

	return db
}

// setup
func initDb() {
	db := connectDb()
	query, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

// example database usage
func getAllMessages(w http.ResponseWriter, r *http.Request) {
	db := connectDb()
	query := "SELECT * FROM message"
	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var messages []Message

	for result.Next() {
		var msg Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Text, &msg.Pubdate, &msg.Flagged)
		if err != nil {
			panic(err.Error())
		}
		w.Write([]byte(msg.Text))
		messages = append(messages, msg)
	}
}

/****************************************
*		  		 PONGO2					*
****************************************/
/*var tplExample = pongo2.Must(pongo2.FromFile("templates/test.html"))

func testPage(w http.ResponseWriter, r *http.Request) {
	err := tplExample.ExecuteWriter(pongo2.Context{"first_name": r.FormValue("deez?")}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}*/

/****************************************
*		   	 GO TEMPLATES				*
****************************************/
var t *template.Template

// findings: variable names MUST be capitalized? idk why
type User struct {
	User bool
	Name string
}

type Data struct {
	FirstName string
	LastName  string
	G         User
}

func initTemplates() {
	t = template.Must(template.ParseFiles("templates/test_layout.html", "templates/test.html"))
}

func url_css() string {
	return "static/style.css"
}

func templatetest(w http.ResponseWriter, r *http.Request) {
	data := Data{
		"Kaare",
		"Børsting",
		User{
			true,
			"boer",
		},
	}
	t.ExecuteTemplate(w, "layout", data)
}

/****************************************
*			 	HANDLERS				*
****************************************/
func hello(w http.ResponseWriter, r *http.Request) {
	data := "Hello World"
	w.Write([]byte(data))
	return
}

/****************************************
*			 	  MAIN					*
****************************************/
func main() {
	initTemplates()

	r := mux.NewRouter()

	r.HandleFunc("/hello", hello)

	// pongo2 template testpage
	// r.HandleFunc("/test", testPage)

	// go templates
	r.HandleFunc("/template", templatetest)

	// Serve static files (css)
	fileServer := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Serve index page on all unhandled routes
	r.PathPrefix("/login").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/login.html")
	})

	r.PathPrefix("/register").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/register.html")
	})

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/timeline.html")
	})

	http.ListenAndServe(":8080", r)
}
