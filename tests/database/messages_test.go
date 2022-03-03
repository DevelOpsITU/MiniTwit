package database

import (
	"fmt"
	"minitwit/config"
	"minitwit/database"
	"testing"
)

var test_user_id = 0

func init() {
	config.SetupTestConfig()
	database.InitGorm()
	test_user_id = database.GormAddUserToDb(test_registration_user)
}

func TestAddMessage(t *testing.T) {
	//TODO: Use test user to post messages

	err := database.AddMessage(test_user_id, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}

	t.Cleanup(func() {
		database.GormRemoveUserFromDb(test_user_id)
	})

}

func TestAddMessageFakeUser(t *testing.T) {

	//TODO: Use a id, that is not of a valid user
	// Right now it does not matter, since we don't use foreign key

	err := database.AddMessage(10000, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}

}

func TestPersonalTimelineMessages(t *testing.T) {

	/*
		This is the correct SQL, which currently returns something else, that we get from GORM:

		select message.message_id,message.author_id,user.username,message.text,message.pub_date, user.email
			from message, user
			where message.flagged = 0 and message.author_id = user.user_id and (
				user.user_id = ? or
			user.user_id in (select whom_id from follower
				where who_id = ?))
				order by message.pub_date desc limit 30


	*/

	messages := database.GetPersonalTimelineMessages(1)

	if len(messages) != 2 {
		t.Errorf("User 1 should have 2 messages to show")
	}

}

func TestUserMessages(t *testing.T) {

	database.AddMessage(test_user_id, "Test message")
	result, err := database.GormGetUserMessages(test_user_id)

	if err != nil {
		t.Errorf(err.Error())
	} else if len(result) != 1 {
		t.Errorf("UserId: " + fmt.Sprint(test_user_id) + " should only have one message, but had " + fmt.Sprint(len(result)))
	}

	if result[0].Text != "Test message" {
		t.Errorf("Message should be 'Test message', but was '" + result[0].Text + "'")
	}

	t.Cleanup(func() {
		database.GormRemoveUserFromDb(test_user_id)
	})

}
