package main

import (
	"fmt"
	"minitwit/config"
	"minitwit/controllers"
	"minitwit/database"

	"gorm.io/driver/postgres"
)

func main() {

	// Get configuration
	config.SetupConfig()

	//database.TestConnection()
	//TODO: Test if the database connection can be established
	// _, err := database.InitGorm(sqlite.Open(config.GetConfig().Database.ConnectionString))
	_, err := database.InitGorm(postgres.Open(config.GetConfig().Database.ConnectionString))
	if err != nil {
		if err != nil {
			fmt.Println(err.Error())
			panic("failed to connect database")
		}
	}
	// Blocking call in router.run
	controllers.HandleRESTRequests()
}
