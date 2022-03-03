package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"minitwit/config"
	"minitwit/controllers"
	"minitwit/database"
)

func main() {

	// Get configuration
	config.SetupConfig()

	database.TestConnection()
	_, err := database.InitGorm(sqlite.Open(config.GetConfig().Database.ConnectionString))
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
