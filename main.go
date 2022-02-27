package main

import (
	"minitwit/config"
	"minitwit/controllers"
	"minitwit/database"
)

func main() {

	database.TestConnection()
	// Get configuration
	config.SetupConfig()
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
