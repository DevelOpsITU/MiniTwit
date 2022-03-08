package database

import (
	"github.com/stretchr/testify/assert"
	"minitwit/database"
	"testing"
)

func init() {

}

func TestFollowUser(t *testing.T) {
	setupTest()
	err := database.FollowUser(user1Id, user2Id)
	if err != nil {
		t.Errorf("User 1 could not follow User 2")
	}

}

func TestFollowUserAgain(t *testing.T) {
	setupTest()
	err := database.FollowUser(user1Id, user2Id)
	if err != nil {
		t.Errorf("User 1 could not follow User 2")
	}

	err = database.FollowUser(user1Id, user2Id)

	assert.NotEmpty(t, err)
	if err == nil {
		t.Errorf("Expected an error that the user is already following that user")
	}

}

func TestFollowUserThatDontExists(t *testing.T) {
	setupTest()
	err := database.FollowUser(user1Id, 9999)
	if err == nil {
		t.Errorf("A user is not allowed to follow a user that does not exist")
	}
}

func TestUnFollowUser(t *testing.T) {
	setupTest()

	following, err := database.GetFollowingUsers(user1Id)
	assert.Empty(t, following)
	assert.Empty(t, err)

	database.FollowUser(user1Id, user2Id)

	following, err = database.GetFollowingUsers(user1Id)
	assert.Equal(t, 1, len(following))

	err = database.UnFollowUser(user1Id, user2Id)
	assert.Empty(t, err)

	following, err = database.GetFollowingUsers(user1Id)
	assert.Empty(t, following)
	assert.Empty(t, err)

}

func TestUnFollowUserAgain(t *testing.T) {
	setupTest()

	following, err := database.GetFollowingUsers(user1Id)
	assert.Empty(t, following)
	assert.Empty(t, err)

	database.FollowUser(user1Id, user2Id)

	following, err = database.GetFollowingUsers(user1Id)
	assert.Equal(t, 1, len(following))
	assert.Empty(t, err)

	err = database.UnFollowUser(user1Id, user2Id)
	assert.Empty(t, err)

	err = database.UnFollowUser(user1Id, user2Id)
	assert.NotEmpty(t, err)

	following, err = database.GetFollowingUsers(user1Id)
	assert.Empty(t, following)
	assert.Empty(t, err)

}
