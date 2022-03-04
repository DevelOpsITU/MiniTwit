package logic

import (
	"errors"
	"minitwit/database"
	"minitwit/models"
)

func AddMessageFromUsername(username string, message string) error {

	if username == "" || message == "" {
		return errors.New("error cant be empty")
	}

	user, err := database.GetUserFromDb(username)

	if err != nil {
		return err
	}

	err = database.AddMessage(user.UserId, message)

	if err != nil {
		return errors.New("message could not be added")
	}
	return nil

}

func AddMessageFromUserModel(user models.User, message string) error {

	if message == "" {
		return errors.New("error cant be empty")
	}

	err := database.AddMessage(user.User_id, message)

	if err != nil {
		return errors.New("message could not be added")
	}
	return nil

}
