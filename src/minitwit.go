package main

import (
	controllers "minitwit/src/controllers"
	"minitwit/src/database"
	config "minitwit/src/settings"
)

func main() {

	database.TestConnection()
	// Get configuration
	config.SetupConfig()
	// Blocking call in router.run
	controllers.HandleRESTRequests()

}
