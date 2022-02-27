package main

import (
	"minitwit/config"
	"minitwit/controllers"
	"minitwit/database"
)

func main() {

	// Get configuration
	config.SetupConfig()
	
	database.TestConnection()
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
