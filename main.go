package main

import (
	"minitwit/controllers"
	"minitwit/database"
	config "minitwit/settings"
)

func main() {

	database.TestConnection()
	// Get configuration
	config.SetupConfig()
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
