package tests

import (
	"database/sql"
	"fmt"
	"minitwit/logic"
	"minitwit/models"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

/*************************************************
* GetUserTwits
**************************************************/
func Test_Get_From_Non_Existing_User_Returns_Empty(t *testing.T) {
	prepare()
	twits, user, err := logic.GetUserTwits("NonExistingUser")
	assert.Empty(t, twits)
	assert.Empty(t, user)
	assert.NotNil(t, err)
}

func Test_Get_From_Existin_User_Returns_Twits(t *testing.T) {
	prepare()
	twits, user, _ := logic.GetUserTwits("Roger Histand")
	twit := twits[0]
	assert.Equal(t, twit.Username, "Roger Histand")
	assert.Equal(t, twit.Pub_date, 1233065594)
	assert.Equal(t, twit.Text, "From hour to hour yesterday I saw my white face of it?")
	assert.Equal(t, user.Username, "Roger Histand")
	assert.Equal(t, user.User_id, 1)
	assert.Equal(t, user.Email, "Roger+Histand@hotmail.com")
}

/*************************************************
* GetPublicTimelineTwits
**************************************************/
func Test_Get_Public_Timline_Returns_Twits(t *testing.T) {
	prepare()
	twits, _ := logic.GetPublicTimelineTwits()
	twit := twits[0]
	assert.Equal(t, twit.Username, "Roger Histand")
	assert.Equal(t, twit.Pub_date, 1233065594)
	assert.Equal(t, twit.Text, "From hour to hour yesterday I saw my white face of it?")
}

/*************************************************
* GetPersonalTimelineTwits
**************************************************/
func Test_Get_Personal_Timline_Returns_Twits(t *testing.T) {
	prepare()
	twits, _ := logic.GetPublicTimelineTwits()
	twit := twits[0]
	assert.Equal(t, twit.Username, "Roger Histand")
	assert.Equal(t, twit.Pub_date, 1233065594)
	assert.Equal(t, twit.Text, "From hour to hour yesterday I saw my white face of it?")
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
func Test_Get_Non_Existing_Gavatar_Should_Fail(t *testing.T) {
	email := "user@mail.com"
	size := "48"
	sizeInt := 48
	hashStr := "6ad193f57f79ac444c3621370da955e9"
	gavatarStr := logic.GetGavateUrl(email, sizeInt)

	expectedString := "http://www.gravatar.com/avatar/" + hashStr + "?d=identicon&s=" + size

	assert.NotEqual(t, gavatarStr, expectedString, "Gavtar failed")
}

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("sqlite3", "minitwit_test.db")
	if err != nil {
		panic("Failed to open db")
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("sqlite"), // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Paths(
			"fixtures",
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
	fmt.Println(fixtures)
	if err := fixtures.Load(); err != nil {
		panic(err.Error())
	}
}

// for code coverage
// https://docs.coveralls.io/go
