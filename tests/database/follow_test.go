package database

import (
	"minitwit/database"
	"testing"
)

func init() {
	print("lol")
}

func TestFollowUser(t *testing.T) {
	setupTest()
	err := database.FollowUser(user1Id, user2Id)
	if err != nil {
		t.Errorf("User 1 could not follow User 2")
	}

}

func TestFollowUserThatDontExists(t *testing.T) {
	setupTest()
	err := database.FollowUser(user1Id, 9999)
	if err == nil {
		t.Errorf("A user is not allowed to follow a user that does not exist")
	}

}
