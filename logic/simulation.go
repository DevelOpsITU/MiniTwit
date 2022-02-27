package logic

import (
	"log"
	"minitwit/src/database"
	"minitwit/src/models"
)

func UnFollowSimulationUser(userId int, user models.User) error {
	err := database.UnFollowUser(userId, user.User_id)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func FollowSimulationUser(userId int, user models.User) error {

	err := database.FollowUser(userId, user.User_id)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetUsernameOfWhoFollowsUser(userId int, noFollowers string) []string {
	users, err := database.GetUsernameOfWhoFollowsUser(userId, noFollowers)

	if err != nil {
		log.Fatal(err)
	}

	return users
}

func GetAllSimulationMessages(noFollowers string) []models.Message {
	users, err := database.GetAllSimulationMessages(noFollowers)

	if err != nil {
		log.Fatal(err)
	}

	return users
}

func GetUserSimulationMessages(user models.User, noMessages string) []models.Message {
	messages, err := database.GetUserSimulationMessages(user.User_id, noMessages)

	if err != nil {
		log.Fatal(err)
	}
	return messages
}
