package main

import (
	"minitwit/src/controllers"
)

func main() {

	// Blocking call in router.run
	controllers.HandleRESTRequests()

}
