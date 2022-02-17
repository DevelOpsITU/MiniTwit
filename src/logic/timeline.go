package logic

import (
	"minitwit/src/database"
	"minitwit/src/models"
)

func GetUserMessages(username string) ([]models.Message, models.User, error) {
	user, err := database.GetUserFromDb(username)

	if err != nil {
		return []models.Message{}, models.User{}, err
	} else {
		return database.GetUserMessages(user.User_id), user, nil
	}

}
