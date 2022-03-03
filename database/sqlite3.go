package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"minitwit/config"
	"minitwit/models"
	"os"
)

const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

/****************************************
*			DATABASE RELATED			*
****************************************/
func ConnectDb() *sql.DB {
	connectionString := config.GetConfig().Database.ConnectionString
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		panic(err)
	}

	return db
}

// setup

func TestConnection() {
	//For Sqlite we simply look for the database file
	connectionString := config.GetConfig().Database.ConnectionString
	if _, err := os.Stat(connectionString); errors.Is(err, os.ErrNotExist) {
		//Does not exist
		fmt.Fprintln(os.Stderr,
			"\n--------------------------------------------------------------\n"+
				"\t File "+connectionString+" does not exists, exiting..\n"+
				"--------------------------------------------------------------")
		os.Exit(1)
	}
}

func UnFollowUser(userId uint, UserIdToUnFollow uint) error {

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

// SIMULATION HANDLING
func HandleSqlQuery(sqlQuery string, args ...interface{}) (*sql.Rows, error) {
	db := ConnectDb()
	query, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	rows, err := query.Query(args...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return rows, nil
}

func GetUsernameOfWhoFollowsUser(userId uint, noFollowers string) ([]string, error) {
	var rows *sql.Rows
	var err error
	if noFollowers == "" {
		query := "SELECT user.username FROM user INNER JOIN follower ON follower.whom_id=user.user_id WHERE follower.who_id=?"
		rows, err = HandleSqlQuery(query, userId)
		if err != nil {
			return []string{}, err
		}
	} else {
		query := "SELECT user.username FROM user INNER JOIN follower ON follower.whom_id=user.user_id WHERE follower.who_id=? LIMIT ?"
		rows, err = HandleSqlQuery(query, userId, noFollowers)
		if err != nil {
			return []string{}, err
		}
	}

	defer rows.Close()
	var userNames []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return []string{}, errors.New("mapping to user error")
		}
		userNames = append(userNames, username)
	}

	return userNames, nil
}

func GetAllSimulationMessages(noFollowers string) ([]models.Message, error) {
	var rows *sql.Rows
	var err error
	if noFollowers == "" {
		query := "SELECT message.text, message.pub_date, user.username from message, user WHERE message.flagged = 0 and message.author_id = user.user_id ORDER BY message.pub_date DESC"
		rows, err = HandleSqlQuery(query)
		if err != nil {
			return []models.Message{}, err
		}
	} else {
		query := "SELECT message.text, message.pub_date, user.username from message, user WHERE message.flagged = 0 and message.author_id = user.user_id ORDER BY message.pub_date DESC LIMIT ?"
		rows, err = HandleSqlQuery(query, noFollowers)
		if err != nil {
			return []models.Message{}, err
		}
	}

	defer rows.Close()
	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Text, &msg.Pubdate, &msg.Username)
		if err != nil {
			return []models.Message{}, errors.New("mapping to user error")
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func GetUserSimulationMessages(userId uint, noFollowers string) ([]models.Message, error) {
	var rows *sql.Rows
	var err error
	if noFollowers == "" {
		query := "SELECT message.text, message.pub_date, user.username FROM message, user WHERE message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ? ORDER BY message.pub_date DESC"
		rows, err = HandleSqlQuery(query, userId)
		if err != nil {
			return []models.Message{}, err
		}
	} else {
		query := "SELECT message.text, message.pub_date, user.username FROM message, user WHERE message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ? ORDER BY message.pub_date DESC LIMIT ?"
		rows, err = HandleSqlQuery(query, userId, noFollowers)
		if err != nil {
			return []models.Message{}, err
		}
	}

	defer rows.Close()
	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Text, &msg.Pubdate, &msg.Username)
		if err != nil {
			return []models.Message{}, errors.New("mapping to user error")
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
