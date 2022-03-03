package database

import (
	"minitwit/database"
	"minitwit/models"
	"testing"
)

func init() {
	_, err := database.InitGorm()
	if err != nil {
		return
	}
}

var test_registration_user = models.RegistrationUser{
	Username:  "a",
	Email:     "@",
	Password1: "1234",
	Password2: "1234",
}

func TestAddUserToDb(t *testing.T) {

	amountOfUsers_old := database.NumberOfUsers()
	database.GormAddUserToDb(test_registration_user)
	amountOfUsers := database.NumberOfUsers()

	if amountOfUsers == amountOfUsers_old {
		t.Error("Users before: ", amountOfUsers_old, " Users now: ", amountOfUsers, " Should be different")
	}

}
