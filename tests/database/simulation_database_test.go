package database

import (
	"github.com/stretchr/testify/assert"
	"minitwit/database"
	"minitwit/models"
	"testing"
)

func init() {

}

//region GormGetUserNameOfWhoUserFollows
func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithNoFollowers_NoNotSpecified_ReturnsEmpty(t *testing.T) {
	setupTest()

	// Arrange
	var usernames []string
	var userId = user4Id

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "")

	// Assert
	assert.Empty(t, usernames)
}

func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithFollowers_NoNotSpecified_ReturnsNotEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var usernames []string
	var userId = user1Id

	database.FollowUser(user1Id, user2Id)

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "")

	// Assert
	assert.NotEmpty(t, usernames)
}

func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithFollowers_NoSetTo0_ReturnsEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var usernames []string
	var userId = user1Id

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "0")

	// Assert
	assert.Empty(t, usernames)
}

func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithFollowers_NoSetTo1_ReturnsOneElement(t *testing.T) {
	setupTest()
	// Arrange
	var usernames []string
	var userId = user1Id
	database.FollowUser(user1Id, user2Id)

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "1")

	// Assert
	assert.NotEmpty(t, usernames)
	assert.Equal(t, 1, len(usernames))

	database.FollowUser(user1Id, user3Id)

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "2")

	// Assert
	assert.NotEmpty(t, usernames)
	assert.Equal(t, 2, len(usernames))
}

//endregion

//region GormGetAllSimulationMessages
func Test_GormGetAllSimulationMessages_NoNotSpecified_ReturnsNotEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var messages []models.Message
	database.AddMessage(user1Id, "My test message")

	// Act
	messages, _ = database.GormGetAllSimulationMessages("")

	// Assert
	assert.NotEmpty(t, messages)
}
func Test_GormGetAllSimulationMessages_NoSetTo1_ReturnsOneElement(t *testing.T) {
	setupTest()
	// Arrange
	var messages []models.Message
	database.AddMessage(user1Id, "My test message")

	// Act
	messages, _ = database.GormGetAllSimulationMessages("1")

	// Assert
	assert.NotEmpty(t, messages)
	assert.Equal(t, 1, len(messages))
}

//endregion

//region GormGetUserSimulationMessages
func Test_GormGetUserSimulationMessages_ExistingUserWithMessages_NoNotSpecified_ReturnsNotEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var messages []models.Message
	var userId = user1Id
	database.AddMessage(userId, "My very own test message")

	// Act
	messages, _ = database.GormGetUserSimulationMessages(userId, "")

	// Assert
	assert.NotEmpty(t, messages)
}

func Test_GormGetUserSimulationMessages_ExistingUserWithMessages__NoSetTo1_ReturnsOneElement(t *testing.T) {
	setupTest()
	// Arrange
	var messages []models.Message
	var userId = user1Id
	database.AddMessage(userId, "My very own test message")

	// Act
	messages, _ = database.GormGetUserSimulationMessages(userId, "1")

	// Assert
	assert.NotEmpty(t, messages)
	assert.Equal(t, 1, len(messages))
}

func Test_GormGetUserSimulationMessages_ExistingUserWithNoMessages_ReturnsEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var messages []models.Message
	var userId = user2Id
	database.AddMessage(user1Id, "Message from user1Id")

	// Act
	messages, _ = database.GormGetUserSimulationMessages(userId, "")

	// Assert
	assert.Empty(t, messages)
}

//endregion
