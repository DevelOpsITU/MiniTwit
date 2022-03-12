package logic

import (
	"minitwit/database"
	"minitwit/log"
	"minitwit/models"
)

func GetAllSimulationMessages(limitStr string) []models.Message {
	users, err := database.GetAllSimulationMessages(limitStr)

	if err != nil {
		log.Logger.Error().Str("limit", limitStr).Msg("Could not get all simulation messages")
	}

	return users
}
