package database

import (
	"gorm.io/driver/sqlite"
	"minitwit/database"
	"minitwit/models"
)

var user1Id uint
var user2Id uint
var user3Id uint
var user4Id uint

var user1 = models.RegistrationUser{
	Username:  "testUser1",
	Email:     "testuser1@mail.com",
	Password1: "pass",
	Password2: "pass",
}

var user2 = models.RegistrationUser{
	Username:  "testUser2",
	Email:     "testuser2@mail.com",
	Password1: "pass",
	Password2: "pass",
}

var user3 = models.RegistrationUser{
	Username:  "testUser3",
	Email:     "testuser3@mail.com",
	Password1: "pass",
	Password2: "pass",
}
var user4 = models.RegistrationUser{
	Username:  "testUser4",
	Email:     "testuser4@mail.com",
	Password1: "pass",
	Password2: "pass",
}

func setupTest() {
	println("Let the testing begin!")

	database.InitGorm(sqlite.Open("file::memory:"))

	user1Id = database.AddUserToDb(user1)
	user2Id = database.AddUserToDb(user2)
	user3Id = database.AddUserToDb(user3)
	user4Id = database.AddUserToDb(user4)

}
