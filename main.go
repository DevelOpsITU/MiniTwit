package main

import (
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

	log.Logger.Info().Msg("Starting MiniTwit application startup checks - complete")
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
