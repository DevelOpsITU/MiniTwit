package database

import (
	"fmt"
	"minitwit/database"
	"testing"
)

func init() {

}

func TestAddMessage(t *testing.T) {
	setupTest()
	err := database.AddMessage(user4Id, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}
}

func TestAddMessageFakeUser(t *testing.T) {
	setupTest()

	// Use a id, that is not of a valid user
	// Right now it does not matter, since we don't use foreign key

	err := database.AddMessage(9999999, "Test message")
	if err != nil {
		t.Errorf("Using a non existing user should have returned an Error!")
	}
}

func TestPersonalTimelineMessages(t *testing.T) {

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

	messages := database.GetPersonalTimelineMessages(user1Id)

	if len(messages) != 1 {
		t.Errorf("User 1 should have its own message 1 but recived: " + fmt.Sprint(len(messages)))
		return
	}

	messages = database.GetPersonalTimelineMessages(user2Id)
	if len(messages) != 2 {
		t.Errorf("User 2 should have the messages from User1 and its own message")
		return
	}

	messages = database.GetPersonalTimelineMessages(user3Id)
	if len(messages) != 3 {
		t.Errorf("User 3 should have the messages from User1, User 2 and its own message")
		return
	}

}

func TestUserMessages(t *testing.T) {

	setupTest()

	database.AddMessage(user1Id, "Test message")
	result, err := database.GetUserMessages(user1Id)

	if err != nil {
		t.Errorf(err.Error())
	} else if len(result) != 1 {
		t.Errorf("UserId: " + fmt.Sprint(user1Id) + " should only have one message, but had " + fmt.Sprint(len(result)))
	}

	if result[0].Text != "Test message" {
		t.Errorf("Message should be 'Test message', but was '" + result[0].Text + "'")
	}

}
