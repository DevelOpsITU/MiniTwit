package database

import (
	"github.com/stretchr/testify/assert"
	"minitwit/database"
	"minitwit/models"
	"testing"
)

func init() {

}

var test_registration_user = models.RegistrationUser{
	Username:  "a",
	Email:     "@",
	Password1: "1234",
	Password2: "1234",
}

func TestAddUserToDb(t *testing.T) {
	setupTest()

	amountOfUsers_old := database.NumberOfUsers()
	user_id := database.AddUserToDb(test_registration_user)
	amountOfUsers := database.NumberOfUsers()

	if amountOfUsers == amountOfUsers_old {
		t.Error("Users before: ", amountOfUsers_old, " Users now: ", amountOfUsers, " Should be different")
	}

	Cleanup_user(user_id)

}

func TestGetUserFromDb(t *testing.T) {
	setupTest()
	test_user_id := database.AddUserToDb(test_registration_user)

	user, err := database.GetUserFromDb(test_registration_user.Username)
	if err != nil {
		t.Error("There should be returned a user")
		return
	}
	assert.Equal(t, test_user_id, user.UserId)
}

func TestCheckIfUserExists(t *testing.T) {
	setupTest()
	database.AddUserToDb(test_registration_user)
	assert.True(t, database.CheckIfUserExists(test_registration_user.Username))
	assert.True(t, database.CheckIfUserExists(user1.Username))
	database.RemoveUserFromDb(user1Id)
	assert.False(t, database.CheckIfUserExists(user1.Username))
	assert.False(t, database.CheckIfUserExists("Non-existing-user"))

}

func Cleanup_user(user_id uint) {
	database.RemoveUserFromDb(user_id)
}
