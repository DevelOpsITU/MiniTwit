package logic

import (
	"errors"
	"minitwit/database"
	"minitwit/log"
	"minitwit/models"
)

func AddMessageFromUsername(username string, message string) error {

	if username == "" || message == "" {
		return errors.New("error cant be empty")
	}

	user, err := database.GetUserFromDb(username)

	if err != nil {
		log.Logger.Error().Str("username", username).Msg("Could not get the user from the database")
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
		log.Logger.Error().Str("username", user.Username).Str("text", message).Msg("Could add the message")
		return errors.New("message could not be added")
	}
	return nil

}
