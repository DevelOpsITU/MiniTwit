package logic

import (
	"errors"
	"minitwit/database"
	"minitwit/log"
	"minitwit/metrics"
	"minitwit/models"
)

func AddMessageFromUsername(username string, message string) error {

	if username == "" || message == "" {
		return errors.New("error cant be empty")
	}

	user, err := database.GetUserFromDb(username)

	if err != nil {
		log.Logger.Warn().Err(err).Caller().Bool("hack", true).Str("username", username).Msg("Created user. Did not exists.")
		user.UserId = database.AddUserToDb(models.RegistrationUser{
			Username:  username,
			Email:     "@",
			Password1: "123",
			Password2: "123",
		})
		metrics.HackCreateUserOnAddMessage.Inc()
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
