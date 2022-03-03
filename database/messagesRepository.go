package database

import (
	"errors"
	"fmt"
	"minitwit/models"
	"strings"
	"time"
)

func GetAllMessages() []models.Message {

	var messages []models.Message

	result, err := gormDb.
		Model(Message{}).
		Order("pub_date desc").
		Limit(30).
		Where("flagged = ?", 0).
		Joins("JOIN user on message.author_id = user.user_id").
		Select("message.message_id , message.author_id , user.username , message.text , message.pub_date , user.email").
		Rows()

	if err != nil {
		panic(err)
	}

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return messages
		}
		messages = append(messages, msg)
	}

	return messages
}

func AddMessage(userId uint, post string) error {

	var message = Message{
		AuthorId:        userId,
		Text:            post,
		PublicationDate: uint(time.Now().Unix()),
		Flagged:         0,
	}

	create := gormDb.Create(&message)

	if create.Error != nil {
		println(create.Error.Error())
		return errors.New(create.Error.Error())
	}

	return nil
}

func GetPersonalTimelineMessages(id uint) []models.Message {

	follows := GetFollowingUsers(id)

	var where string

	if follows == nil {
		where = "flagged = ? AND message.author_id = user.user_id and ( user.user_id = ? )"
	} else {
		where = "flagged = ? AND message.author_id = user.user_id and ( user.user_id = ? or user.user_id in (" + arrayToString(follows, ",") + "))"
	}

	result, err := gormDb.
		Model(&Message{}).
		Joins("JOIN user on message.author_id = user.user_id").
		Order("pub_date desc").
		Limit(30).
		Where(where, 0, id).
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

func arrayToString(a []uint, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func GetUserMessages(userId uint) ([]models.Message, error) {

	var messages []models.Message

	result, err := gormDb.
		Model(Message{}).
		Limit(30).
		Order("pub_date desc").
		Where("message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ?", userId).
		Joins("JOIN user on message.author_id = user.user_id").
		Select("message.message_id, message.author_id, user.username, message.text, message.pub_date, user.email").
		Rows()

	if err != nil {
		return messages, errors.New("Failed to get the userMessages: " + err.Error())
	}

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return messages, errors.New("Failed to scan the element: " + err.Error())
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
