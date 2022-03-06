package database

import (
	"errors"
	"fmt"
	"minitwit/log"
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
		Joins("JOIN \"user\" AS u ON message.author_id = u.user_id").
		Select("message.message_id, message.author_id, u.username, message.text, message.pub_date, u.email").
		Rows()

	if err != nil {
		log.Logger.Error().Err(err).Caller().Msg("Could not get public messages")
		return messages
	}

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			log.Logger.Error().Err(err).Caller().Msg("Could not map the messages")
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
		log.Logger.Error().Str("userId", fmt.Sprint(userId)).Str("text", post).Msg("Could not create the message")
		return errors.New(create.Error.Error())
	}

	return nil
}

func GetPersonalTimelineMessages(id uint) []models.Message {

	var messages []models.Message

	follows, err := GetFollowingUsers(id)

	if err != nil {
		log.Logger.Error().Err(err).Str("userId", fmt.Sprint(id)).Msg("Error getting the users following")
		return messages
	}

	var where string

	if follows == nil {
		where = "flagged = ? AND message.author_id = u.user_id and ( u.user_id = ? )"
	} else {
		where = "flagged = ? AND message.author_id = u.user_id and ( u.user_id = ? or u.user_id in (" + arrayToString(follows, ",") + "))"
	}

	result, err := gormDb.
		Model(&Message{}).
		Joins("JOIN \"user\" AS u ON message.author_id = u.user_id").
		Order("pub_date desc").
		Limit(30).
		Where(where, 0, id).
		Select("message.message_id , message.author_id , u.username , message.text , message.pub_date , u.email").
		Rows()

	if err != nil {
		log.Logger.Error().Str("userId", fmt.Sprint(id)).Msg("Could not get the messages from the user and followers")
		return messages
	}

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			log.Logger.Error().Err(err).Str("userId", fmt.Sprint(id)).Msg("Could not map the message from user")
			return messages
		}
		messages = append(messages, msg)
	}

	return messages
}

func arrayToString(a []uint, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func GetUserMessages(userId uint, limit int) ([]models.Message, error) {

	var messages []models.Message

	result, err := gormDb.
		Model(Message{}).
		Limit(limit).
		Order("pub_date desc").
		Where("message.flagged = 0 AND u.user_id = message.author_id AND u.user_id = ?", userId).
		Joins("JOIN \"user\" AS u ON message.author_id = u.user_id").
		Select("message.message_id, message.author_id, u.username, message.text, message.pub_date, u.email").
		Rows()

	if err != nil {
		log.Logger.Error().Str("userId", fmt.Sprint(userId)).Msg("Could not get the messages from the user")
		return messages, errors.New(err.Error())
	}

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			log.Logger.Error().Str("userId", fmt.Sprint(userId)).Msg("Could not map the message from user")
			return messages, errors.New("Failed to scan the element: " + err.Error())
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
