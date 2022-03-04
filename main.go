package main

import (
	"gorm.io/driver/sqlite"
	"minitwit/config"
	"minitwit/controllers"
	"minitwit/database"
	"minitwit/log"
	"os"
)

func main() {

	// Get configuration
	config.SetupConfig()
	log.SetUpLogger()

	log.Logger.Info().Msg("Starting MiniTwit application startup checks")

	//database.TestConnection()
	//TODO: Test if the database connection can be established
	_, err := database.InitGorm(sqlite.Open(config.GetConfig().Database.ConnectionString))
	if err != nil {
		log.Logger.Error().Err(err).Msg("Unable to connect with the database")
		os.Exit(1)

	}

	log.Logger.Info().Msg("Starting MiniTwit application startup checks - complete")
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
