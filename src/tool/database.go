package tool

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
)

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

type User struct {
	User_id  int
	Username string
	Email    string
	pw_hash  string
}

const DATABASE = "./minitwit.db"

//const DATABASE = "C:/Users/hardk/source/repos/MiniTwit/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

/****************************************
*			DATABASE RELATED			*
****************************************/
func ConnectDb() *sql.DB {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}

	return db
}

// setup
func InitDb() {
	db := ConnectDb()
	query, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

// example Database usage
func GetAllMessages() {
	db := ConnectDb()
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
		messages = append(messages, msg)
	}
	//fmt.Printf("%+v\n", messages)
}

func GetUserFromDb(username string) {
	db := ConnectDb()
	stmt, err := db.Prepare("SELECT x.* FROM 'user' x WHERE username like '?'")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(username)

	var user User
	err = result.Scan(&user.User_id, &user.Username, &user.Email, &user.pw_hash)

	if err != nil {
		log.Fatal(err)
	} else {

	}

}
