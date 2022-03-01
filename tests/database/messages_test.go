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

	t.Error("")

	if len(messages) != 30 {
		t.Errorf("User 1 should have at least 30 messages to show")
	}

}
