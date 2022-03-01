package database

import (
	"minitwit/database"
	"testing"
)

func init() {
	database.InitGorm()
}

func TestAddMessage(t *testing.T) {
	//TODO: Use test user to post messages

	err := database.AddMessage(1, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}

}

func TestAddMessageFakeUser(t *testing.T) {

	//TODO: Use a id, that is not of a valid user
	// Right now it does not matter, since we don't use foreign key

	err := database.AddMessage(10000, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}

}
