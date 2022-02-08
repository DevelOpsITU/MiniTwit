package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//const DATABASE = "tmp/minitwit.db"
const DATABASE = "C:/Users/hardk/source/repos/MiniTwit/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

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
	fmt.Println(query)
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

// example database usage
func getAllMessages() {
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
		fmt.Println(msg.Text)
		messages = append(messages, msg)
	}
}

/****************************************
*			 	HANDLERS				*
****************************************/
func deez(w http.ResponseWriter, r *http.Request) {
	data := "Hello World"
	w.Write([]byte(data))
	return
}

/****************************************
*			 	  MAIN					*
****************************************/
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/deez", deez)
	http.Handle("/", r)
	getAllMessages()
	http.ListenAndServe(":8080", r)
}
