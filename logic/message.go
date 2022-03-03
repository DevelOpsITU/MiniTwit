package logic

import (
	"errors"
	"minitwit/database"
	"minitwit/models"
)

func AddMessage(user models.User, message string) error {

	if message == "" {
		return errors.New("error cant be empty")
	}

	err := database.AddMessage(user.User_id, message)

	if err != nil {
		return errors.New("message could not be added")
	}
	return nil

}
