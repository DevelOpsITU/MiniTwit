package database

import (
	"fmt"
	"minitwit/database"
	"testing"
)

func init() {
	print("lol")
}

func TestAddMessage(t *testing.T) {
	setupTest()
	//TODO: Use test user to post messages
	err := database.AddMessage(user4Id, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}

	database.GormRemoveMessagesFromDb(user4Id)

}

func TestAddMessageFakeUser(t *testing.T) {
	setupTest()

	//TODO: Use a id, that is not of a valid user
	// Right now it does not matter, since we don't use foreign key

	err := database.AddMessage(9999999, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}
	database.GormRemoveMessagesFromDb(9999999)

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

	setupTest()

	database.AddMessage(user1Id, "Message from user1")
	database.AddMessage(user2Id, "Message from user2")
	database.AddMessage(user3Id, "Message from user3")

	// User1 follows no one
	// User2 follows User1
	database.FollowUser(user2Id, user1Id)
	// User3 follows User1 and User2
	database.FollowUser(user3Id, user1Id)
	database.FollowUser(user3Id, user2Id)

	messages := database.GetPersonalTimelineMessages(1)

	if len(messages) != 1 {
		t.Errorf("User 1 should only have its own message")
	}

	messages = database.GetPersonalTimelineMessages(2)
	if len(messages) != 2 {
		t.Errorf("User 2 should have the messages from User1 and its own message")
	}

	messages = database.GetPersonalTimelineMessages(3)
	if len(messages) != 3 {
		t.Errorf("User 2 should have the messages from User1, User 2 and its own message")
	}

}

func TestUserMessages(t *testing.T) {

	setupTest()

	database.AddMessage(user1Id, "Test message")
	result, err := database.GormGetUserMessages(user1Id)

	if err != nil {
		t.Errorf(err.Error())
	} else if len(result) != 1 {
		t.Errorf("UserId: " + fmt.Sprint(user1Id) + " should only have one message, but had " + fmt.Sprint(len(result)))
	}

	if result[0].Text != "Test message" {
		t.Errorf("Message should be 'Test message', but was '" + result[0].Text + "'")
	}
	database.GormRemoveMessagesFromDb(user1Id)

	t.Cleanup(func() {
		database.GormRemoveUserFromDb(user1Id)
	})

}
