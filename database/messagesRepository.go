package database

import (
	"errors"
	"fmt"
	"minitwit/models"
	"strings"
	"time"
)

func GormGetAllMessages() []models.Message {

	result, err := gormDb.
		Model(models.Message{}).
		Table("message").
		Order("pub_date desc").
		Limit(30).
		Where("flagged = ?", 0).
		Joins("JOIN user on message.author_id = user.user_id").
		Select("message.message_id , message.author_id , user.username , message.text , message.pub_date , user.email").
		Rows()

	if err != nil {
		panic(err)
	}

	var messages2 []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages2 = append(messages2, msg)
	}

	return messages2
}

func AddMessage(userId uint, message string) error {

	var messageObj = models.Message{
		AuthorId: userId,
		Text:     message,
		Pubdate:  time.Now().Unix(),
		Flagged:  0,
	}

	create := gormDb.
		Select("message_id", "author_id", "text", "pub_date", "flagged").
		Table("message").Create(&messageObj)

	if create.Error != nil {
		println(create.Error.Error())
		return errors.New(create.Error.Error())
	}

	return nil
}

func GormRemoveMessagesFromDb(user_id uint) {

	result := gormDb.
		Where("author_id = ?", user_id).
		Delete(&Message{})

	if result.Error != nil {
		panic(result.Error)
	}

}

// Returns a list of all the users a user is following
func getFollowingUsers(userId uint) []int {

	var follows []int

	subquery, err := gormDb.
		Model(&Follower{}).
		Select("whom_id").
		Where("who_id = ?", userId).
		Rows()

	if err != nil {
		//TODO: Remove panic statements. it crashes the application.
		panic(err)
	}

	for subquery.Next() {
		var user int
		err := subquery.Scan(&user)
		if err != nil {
			//TODO
		}
		follows = append(follows, user)
	}
	return follows
}

func GetPersonalTimelineMessages(id uint) []models.Message {

	follows := getFollowingUsers(id)
	result, err := gormDb.
		Model(&Message{}).
		Joins("JOIN user on message.author_id = user.user_id").
		Order("pub_date desc").
		Limit(30).
		Where("flagged = ? AND message.author_id = user.user_id and ( user.user_id = ? or user.user_id in (?))",
			0, id, arrayToString(follows, ",")).
		Select("message.message_id , message.author_id , user.username , message.text , message.pub_date , user.email").
		Rows()

	if err != nil {
		//TODO: Remove panic statements. it crashes the application.
		panic(err)
	}

	var messages2 []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages2 = append(messages2, msg)
	}

	return messages2
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func GormGetUserMessages(userId uint) ([]models.Message, error) {
	result, err := gormDb.
		Model(models.Message{}).
		Limit(30).
		Table("message").
		Order("pub_date desc").
		Where("message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ?", userId).
		Joins("JOIN user on message.author_id = user.user_id").
		Select("message.message_id, message.author_id, user.username, message.text, message.pub_date, user.email").
		Rows()

	if err != nil {
		return []models.Message{}, errors.New("Failed to get the userMessages: " + err.Error())
	}

	var messages []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}, errors.New("Failed to scan the element: " + err.Error())
		}
		messages = append(messages, msg)
	}
	defer result.Close()
	return messages, nil
}
