package logic

import (
	"errors"
	"minitwit/src/database"
	"minitwit/src/models"
)

func AddMessage(user models.User, message string) error {

	if message == "" {
		return errors.New("Error cant be empty.")
	}

	err := database.AddMessage(user.User_id, message)

	if err != nil {
		return errors.New("Message could not be added")
	}
	return nil

}
