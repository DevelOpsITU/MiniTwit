package main

import (
	"minitwit/config"
	"fmt"
	"minitwit/controllers"
	"minitwit/database"
)

func main() {

	// Get configuration
	config.SetupConfig()

	database.TestConnection()
	_, err := database.InitGorm()
	if err != nil {
		if err != nil {
			fmt.Println(err.Error())
			panic("failed to connect database")
		}
	}
	database.GormGetAllMessages()

	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
