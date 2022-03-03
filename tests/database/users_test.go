package database

import (
	"gorm.io/driver/sqlite"
	"minitwit/database"
	"minitwit/models"
	"testing"
)

func init() {
	_, err := database.InitGorm(sqlite.Open("file::memory:?cache=shared"))
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
	user_id := database.GormAddUserToDb(test_registration_user)
	amountOfUsers := database.NumberOfUsers()

	if amountOfUsers == amountOfUsers_old {
		t.Error("Users before: ", amountOfUsers_old, " Users now: ", amountOfUsers, " Should be different")
	}

	Cleanup_user(user_id)

}

func Cleanup_user(user_id uint) {
	database.GormRemoveUserFromDb(user_id)
}
