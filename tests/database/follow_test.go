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
