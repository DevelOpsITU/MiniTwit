package main

import (
	backgroundservices "minitwit/background-services"
	"minitwit/config"
	"minitwit/controllers"
	"minitwit/database"
	"minitwit/log"
)

func main() {

	// Get configuration
	config.SetupConfig()
	log.SetUpLogger()

	log.Logger.Info().Msg("Starting MiniTwit application startup checks")

	database.Init()

	backgroundservices.Init()

	log.Logger.Info().Msg("Starting MiniTwit application startup checks - complete")
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
