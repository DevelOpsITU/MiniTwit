package logic

import (
	"log"
	"minitwit/src/database"
	"minitwit/src/models"
)

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
