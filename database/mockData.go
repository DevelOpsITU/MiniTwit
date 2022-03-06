package database

import (
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

func GenerateMockData() {
	println("Let the testing begin!")

	user1Id = AddUserToDb(user1)
	user2Id = AddUserToDb(user2)
	user3Id = AddUserToDb(user3)
	user4Id = AddUserToDb(user4)

	AddMessage(user1Id, "Message from user1")
	AddMessage(user2Id, "Message from user2")
	AddMessage(user3Id, "Message from user3")
	AddMessage(user4Id, "Message from user4")

}
