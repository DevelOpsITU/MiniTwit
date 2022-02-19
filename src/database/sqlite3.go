package database

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"minitwit/src/models"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
)

const DATABASE = "/tmp/minitwit.db"

// const DATABASE = "C:/Users/hardk/source/repos/MiniTwit/minitwit.db"
//const DATABASE = "/home/turbo/ITU/DevOps/MiniTwit/tmp/minitwit.db"
//const DATABASE = "C:\\Users\\JTT\\Documents\\git\\MiniTwit\\minitwit.db"

const DATABASE = "H:/repos/MiniTwit/minitwit.db"

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

func TestConnection() {
	//For Sqlite we simply look for the database file
	if _, err := os.Stat(DATABASE); errors.Is(err, os.ErrNotExist) {
		//Does not exist
		fmt.Fprintln(os.Stderr,
			"\n--------------------------------------------------------------\n"+
				"\t File "+DATABASE+" does not exists, exiting..\n"+
				"--------------------------------------------------------------")
		os.Exit(1)
	}
}

// example Database usage
func GetUserMessages(id int) []models.Message {
	db := ConnectDb()
	query := string(`SELECT 
		message.message_id, 
		message.author_id, 
		user.username, 
		message.text, 
		message.pub_date, 
		user.email 
		FROM message, user 
		WHERE message.flagged = 0 AND 
		user.user_id = (?) AND
		user.user_id = message.author_id
		ORDER BY message.pub_date DESC 
		LIMIT 30`)
	result, err := db.Query(query, fmt.Sprint(id), fmt.Sprint(id))
	if err != nil {
		panic(err)
	}

	var messages []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages = append(messages, msg)
	}
	defer result.Close()
	return messages
}

func GetAllMessages() []models.Message {
	db := ConnectDb()
	query := string("select message.message_id , message.author_id , user.username , message.text , message.pub_date ,  user.email from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit 30")
	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	var messages []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages = append(messages, msg)
	}
	defer result.Close()
	return messages
}

func GetPersonalTimelineMessages(id int) []models.Message {
	db := ConnectDb()
	query := string(`
		select message.message_id,message.author_id,user.username,message.text,message.pub_date, user.email
		from message, user
        where message.flagged = 0 and message.author_id = user.user_id and (
            user.user_id = ? or
            user.user_id in (select whom_id from follower
                                    where who_id = ?))
        order by message.pub_date desc limit 30`)
	result, err := db.Query(query, fmt.Sprint(id), fmt.Sprint(id))
	if err != nil {
		panic(err)
	}

	var messages []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages = append(messages, msg)
	}
	defer result.Close()
	return messages
}

// TODO: Return errors if any, and meybe the user
func AddUserToDb(user models.RegistrationUser) {
	db := ConnectDb()
	salt := make([]byte, 4)
	io.ReadFull(rand.Reader, salt)

	pwIteration_int, _ := strconv.Atoi("50000")
	dk := pbkdf2.Key([]byte(user.Password1), salt, pwIteration_int, 32, sha256.New)

	pw_hashed := "pbkdf2:sha256:50000$" + string(salt) + "$" + hex.EncodeToString(dk)
	query, err := db.Prepare("INSERT INTO user(username, email, pw_hash) values (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(user.Username, user.Email, pw_hashed)

	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()
}

func GetUserFromDb(username string) (models.User, error) {
	db := ConnectDb()
	//TODO: Prepared statements
	strs := []string{"SELECT x.* FROM 'user' x WHERE username like '", username, "'"}
	query := strings.Join(strs, "")
	row, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return models.User{}, errors.New("database connection error")
	}
	var user models.User
	for row.Next() { // Iterate and fetch the records from result cursor

		err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)
		if err != nil {
			return models.User{}, errors.New("Mapping to user error")
		}
	}
	defer row.Close()

	// This is a quazzy hackz when no user is returned
	if user.User_id == 0 {
		return models.User{}, errors.New("User not found")
	}
	return user, nil

}

func CheckIfUserExists(username string) bool {

	user, _ := GetUserFromDb(username)
	if user.User_id != 0 {
		return true
	}

	return false
}

func FollowUser(userId int, UserIdToFollow int) error {
	db := ConnectDb()
	query, err := db.Prepare("INSERT INTO follower (who_id, whom_id) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = query.Exec(userId, UserIdToFollow)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer query.Close()
	return nil
}

func UnFollowUser(userId int, UserIdToUnFollow int) error {

	db := ConnectDb()
	query, err := db.Prepare("DELETE FROM follower WHERE who_id = ? AND whom_id = ?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = query.Exec(userId, UserIdToUnFollow)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer query.Close()
	return nil
}

func AddMessage(userId int, message string) error {

	db := ConnectDb()

	query, err := db.Prepare(`INSERT INTO message (author_id, text, pub_date, flagged) 
			VALUES (?, ?, ?, 0)`)
	if err != nil {
		println(err.Error())
		return err

	}
	_, err = query.Exec(userId, message, time.Now().Unix())

	if err != nil {
		println(err.Error())
		return err
	}
	defer query.Close()

	return nil
}
