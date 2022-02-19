package main

import (
	controllers "minitwit/src/controllers"
	"minitwit/src/database"
)

func main() {

	database.TestConnection()
	// Blocking call in router.run
	controllers.HandleRESTRequests()

}
