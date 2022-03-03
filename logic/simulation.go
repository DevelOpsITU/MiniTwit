package logic

import (
	"log"
	"minitwit/database"
	"minitwit/models"
)

func UnFollowSimulationUser(userId uint, user models.User) error {
	err := database.UnFollowUser(userId, user.User_id)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func FollowSimulationUser(userId uint, user models.User) error {

	err := database.FollowUser(userId, user.User_id)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetUsernameOfWhoFollowsUser(userId uint, noFollowers string) []string {
	users, err := database.GormGetUsernameOfWhoFollowsUser(userId, noFollowers)

	if err != nil {
		log.Fatal(err)
	}

	return users
}

func GetAllSimulationMessages(noFollowers string) []models.Message {
	users, err := database.GormGetAllSimulationMessages(noFollowers)

	if err != nil {
		log.Fatal(err)
	}

	return users
}

func GetUserSimulationMessages(user models.User, noMessages string) []models.Message {
	messages, err := database.GormGetUserSimulationMessages(user.User_id, noMessages)

	if err != nil {
		log.Fatal(err)
	}
	return messages
}
