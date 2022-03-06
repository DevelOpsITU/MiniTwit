package tests

import (
	"minitwit/database"
	"minitwit/logic"
	"minitwit/models"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	fixtures *testfixtures.Loader
)

/*************************************************
* GetUserTwits
**************************************************/
func Test_Get_From_Non_Existing_User_Returns_Empty(t *testing.T) {
	prepare()
	twits, user, err := logic.GetUserTwits("NonExistingUser", 48)
	assert.Empty(t, twits)
	assert.Empty(t, user)
	assert.NotNil(t, err)
}

func Test_Get_From_Existin_User_Returns_Twits(t *testing.T) {
	prepare()
	twits, user, _ := logic.GetUserTwits("Roger Histand", 48)
	twit := twits[0]
	assert.Equal(t, "Roger Histand", twit.Username)
	assert.Equal(t, logic.FormatPubdate(1233065594), twit.Pub_date)
	assert.Equal(t, "From hour to hour yesterday I saw my white face of it?", twit.Text)
	assert.Equal(t, "Roger Histand", user.Username)
	assert.Equal(t, uint(1), user.User_id)
	// assert.Equal(t, user.Email, "Roger+Histand@hotmail.com")
}

/*************************************************
* GetPublicTimelineTwits
**************************************************/
func Test_Get_Public_Timline_Returns_Twits(t *testing.T) {
	prepare()
	twits, _ := logic.GetPublicTimelineTwits()
	twit := twits[0]
	assert.Equal(t, "Roger Histand", twit.Username)
	assert.Equal(t, logic.FormatPubdate(1233065594), twit.Pub_date)
	assert.Equal(t, "From hour to hour yesterday I saw my white face of it?", twit.Text)
}

/*************************************************
* GetPersonalTimelineTwits
**************************************************/
func Test_Get_Personal_Timline_Returns_Twits(t *testing.T) {
	prepare()
	twits, _ := logic.GetPersonalTimelineTwits(models.User{Username: "Roger Histand"})
	twit := twits[0]
	assert.Equal(t, "Roger Histand", twit.Username)
	assert.Equal(t, logic.FormatPubdate(1233065594), twit.Pub_date)
	assert.Equal(t, "From hour to hour yesterday I saw my white face of it?", twit.Text)
}

/*************************************************
* ConvertMessagesToTwits
**************************************************/
func Test_Convert_Empty_Returns_Empty(t *testing.T) {
	prepare()
	messages := make([]models.Message, 0)
	twits := logic.ConvertMessagesToTwits(&messages)
	assert.Empty(t, twits)
}

/*************************************************
* GetGavatarUrl
**************************************************/
func Test_Get_Existing_Gavatar_Should_Return(t *testing.T) {
	email := "user@mail.com"
	size := "48"
	sizeInt := 48
	hashStr := "6ad193f57f79ac444c3621370da955e9"
	gavatarStr := logic.GetGavaterUrl(email, sizeInt)

	expectedString := "http://www.gravatar.com/avatar/" + hashStr + "?d=identicon&s=" + size

	assert.Equal(t, expectedString, gavatarStr)
}

func TestMain(m *testing.M) {
	var err error
	db, err = database.InitGorm(sqlite.Open("file::memory:"))

	if err != nil {
		panic("Failed to init gorm")
	}

	data, err := db.DB()

	if err != nil {
		panic("Failed to fetch sql db")
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(data),
		testfixtures.Dialect("sqlite"), // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Paths(
			"../../fixtures",
		), // YAML files
	)
	if err != nil {
		panic("Failed to create fixtures")
	}

	// run tests
	exitVal := m.Run()

	os.Exit(exitVal)
}

func prepare() {
	if err := fixtures.Load(); err != nil {
		panic(err.Error())
	}
}

// for code coverage
// https://docs.coveralls.io/go
