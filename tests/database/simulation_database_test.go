package database

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"minitwit/database"
	"minitwit/models"
	"testing"
)

func init() {
	database.InitGorm(sqlite.Open("file::memory:?cache=shared"))
}

//region GormGetUserNameOfWhoUserFollows
func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithNoFollowers_NoNotSpecified_ReturnsEmpty(t *testing.T) {
	setupTest()

	// Arrange
	var usernames []string
	var userId = 206 // TODO: Edit this to a test-user that does not follow anyone

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "")

	// Assert
	assert.Empty(t, usernames)
}

func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithFollowers_NoNotSpecified_ReturnsNotEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var usernames []string
	var userId = 1 // TODO: Edit this to a test-user that does follow anyone

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "")

	// Assert
	assert.NotEmpty(t, usernames)
}

func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithFollowers_NoSetTo0_ReturnsEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var usernames []string
	var userId = 1 // TODO: Edit this to a test-user that does follow anyone

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "0")

	// Assert
	assert.Empty(t, usernames)
}

func Test_GormGetUserNameOfWhoUserFollows_ExistingUserWithFollowers_NoSetTo1_ReturnsOneElement(t *testing.T) {
	setupTest()
	// Arrange
	var usernames []string
	var userId = 1 // TODO: Edit this to a test-user that does follow anyone

	// Act
	usernames, _ = database.GormGetUsernameOfWhoFollowsUser(userId, "1")

	// Assert
	assert.NotEmpty(t, usernames)
	assert.Equal(t, 1, len(usernames))
}

//endregion

//region GormGetAllSimulationMessages
func Test_GormGetAllSimulationMessages_NoNotSpecified_ReturnsNotEmpty(t *testing.T) {
	setupTest()
	// Arrange
	var messages []models.Message
	// TODO: Ensure there is at least one message in database

	// Act
	messages, _ = database.GormGetAllSimulationMessages("")

	// Assert
	assert.NotEmpty(t, messages)
}
func Test_GormGetAllSimulationMessages_NoSetTo1_ReturnsOneElement(t *testing.T) {
	setupTest()
	// Arrange
	var messages []models.Message
	// TODO: Ensure there is at least one message in database

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
	var userId = user2Id // TODO: Edit this to a test-user that has messages
	// TODO: Ensure there is at least one message in database

	// Act
	messages, _ = database.GormGetUserSimulationMessages(userId, "")

	// Assert
	assert.Empty(t, messages)
}

//endregion
