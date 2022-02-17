package logic

import (
	"errors"
	"minitwit/src/database"
	"minitwit/src/models"
	"strings"
)

func CreateUser(registationUser models.RegistrationUser) (models.User, error) {

	if registationUser.Username == "" {
		return models.User{}, errors.New("You have to enter a username")
	} else if registationUser.Email == "" || !strings.Contains(registationUser.Email, "@") {
		return models.User{}, errors.New("Your have to enter a valid email address")
	} else if registationUser.Password1 == "" {
		return models.User{}, errors.New("You have to enter a password")
	} else if registationUser.Password1 != registationUser.Password2 {
		return models.User{}, errors.New("The two passwords do not match")
	} else {

		user, err := database.GetUserFromDb(registationUser.Username)
		if err != nil {
			return models.User{}, errors.New("The username is already taken")
		} else {
			database.AddUserToDb(registationUser)
			return user, nil
		}
	}
}
