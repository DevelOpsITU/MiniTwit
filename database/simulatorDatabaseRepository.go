package database

import (
	"database/sql"
	"errors"
	"minitwit/models"
	"strconv"
)

// TODO: This name of method is probably not the best
func GetUsernameOfWhoFollowsUser(userId uint, noFollowers string) ([]string, error) {
	var result *sql.Rows
	var err error
	if noFollowers == "" {
		result, err = gormDb.
			Model(models.User{}).
			Table("user").
			Where("follower.who_id=?", userId).
			Joins("INNER JOIN follower ON follower.whom_id=user.user_id").
			Select("user.username").
			Rows()
	} else {
		var limit int
		limit, err = strconv.Atoi(noFollowers)
		if limit == 0 {
			return []string{}, nil
		}
		result, err = gormDb.
			Model(models.User{}).
			Table("user").
			Limit(limit).
			Where("follower.who_id=?", userId).
			Joins("INNER JOIN follower ON follower.whom_id=user.user_id").
			Select("user.username").
			Rows()
	}

	if err != nil {
		return []string{}, err
	}

	var usernames []string
	for result.Next() {
		var username string
		err := result.Scan(&username)
		if err != nil {
			return []string{}, errors.New("mapping to user error")
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}

func GetAllSimulationMessages(limitStr string) ([]models.Message, error) {
	var result *sql.Rows
	var err error
	if limitStr == "" {
		result, err = gormDb.
			Model(models.User{}).
			Table("message").
			Order("pub_date desc").
			Where("message.flagged = 0 and message.author_id = user.user_id").
			Joins("JOIN user on message.author_id = user.user_id").
			Select("message.text, message.pub_date, user.username").
			Rows()
	} else {
		var limit int
		limit, err = strconv.Atoi(limitStr)
		if limit == 0 {
			return []models.Message{}, nil
		}
		result, err = gormDb.
			Model(models.User{}).
			Table("message").
			Limit(limit).
			Where("message.flagged = 0 and message.author_id = user.user_id").
			Joins("JOIN user on message.author_id = user.user_id").
			Order("pub_date desc").
			Select("message.text, message.pub_date, user.username").
			Rows()
	}

	if err != nil {
		return []models.Message{}, err
	}

	var messages []models.Message
	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.Text, &msg.Pubdate, &msg.Username)
		if err != nil {
			return []models.Message{}, errors.New("mapping to message error")
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func GetUserSimulationMessages(userId uint, limitStr string) ([]models.Message, error) {
	var result *sql.Rows
	var err error
	if limitStr == "" {
		result, err = gormDb.
			Model(models.User{}).
			Table("message").
			Order("pub_date desc").
			Where("message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ?", userId).
			Joins("JOIN user on message.author_id = user.user_id").
			Select("message.text, message.pub_date, user.username").
			Rows()
	} else {
		var limit int
		limit, err = strconv.Atoi(limitStr)
		if limit == 0 {
			return []models.Message{}, nil
		}
		result, err = gormDb.
			Model(models.User{}).
			Table("message").
			Limit(limit).
			Order("pub_date desc").
			Where("message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ?", userId).
			Joins("JOIN user on message.author_id = user.user_id").
			Select("message.text, message.pub_date, user.username").
			Rows()
	}

	if err != nil {
		return []models.Message{}, err
	}

	var messages []models.Message
	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.Text, &msg.Pubdate, &msg.Username)
		if err != nil {
			return []models.Message{}, errors.New("mapping to message error")
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
