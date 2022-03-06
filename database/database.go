package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"minitwit/config"
	"minitwit/log"
	"os"
)

func Init() {

	var err error
	if config.GetConfig().Database.Type == "SQLite" {
		_, err = InitGorm(sqlite.Open(config.GetConfig().Database.ConnectionString))

	} else if config.GetConfig().Database.Type == "Postgres" {
		_, err = InitGorm(postgres.Open(config.GetConfig().Database.ConnectionString))
	} else {
		err = errors.New("database type not selected")
	}
	if err != nil {
		if err != nil {
			log.Logger.Err(err).Msg("Could not open database")
			os.Exit(1)
		}
	}

	if config.GetConfig().Development.GenerateMockData == true {
		if config.GetConfig().Database.ConnectionString == "file::memory:" && config.GetConfig().Database.Type == "SQLite" {
			GenerateMockData()
		}
	}

}
