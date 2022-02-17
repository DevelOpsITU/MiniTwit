package logic

import (
	"errors"
	"log"
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

func FollowUser(userId int, usernameToFollow string) error {

	userToFollow, err := database.GetUserFromDb(usernameToFollow)

	if err != nil {
		return err
	}

	//TODO: Check that the user is not already following the user Issue #47
	err = database.FollowUser(userId, userToFollow.User_id)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func UnFollowUser(userId int, usernameToUnFollow string) error {
	userToUnFollow, err := database.GetUserFromDb(usernameToUnFollow)

	if err != nil {
		return err
	}
	// TODO: check if already following before trying this
	err = database.UnFollowUser(userId, userToUnFollow.User_id)

	if err != nil {
		log.Fatal(err)
	}

	return nil

}
