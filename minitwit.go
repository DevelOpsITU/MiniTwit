package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const DATABASE = "tmp/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

/****************************************
*			DATABASE RELATED			*
****************************************/
func connectDb() *sql.DB {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
	val, err := db.Query("SELECT table_name FROM all_tables;")
	fmt.Println(val)
}

/*func queryDb(query, args, one) {
	cur :=
}*/

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
	initDb()
	http.ListenAndServe(":8080", r)
}
