package main

import (
	"minitwit/controllers"
	"minitwit/database"
)

func main() {

	database.TestConnection()
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
