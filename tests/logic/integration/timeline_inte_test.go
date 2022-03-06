package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"minitwit/controllers"
	"minitwit/database"

	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	fixtures *testfixtures.Loader
)

func setup() *gin.Engine {
	router := gin.Default()

	for _, handler := range *controllers.GetHandlers() {
		handler.(func(engine *gin.Engine))(router)
	}

	return router
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	return w
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

/*************************************************
* GetPublicTimelineTwits
**************************************************/
func Test_Get_Public_Timline_Returns_Twits(t *testing.T) {
	prepare()
	test_server := httptest.NewServer(setup())
	defer test_server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/public", test_server.URL))
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)
}

/*************************************************
* GetUserTwits
**************************************************/
func Test_Get_From_Non_Existing_User_Returns_Empty(t *testing.T) {
	// expected result
	body := gin.H{}

	prepare()
	router := setup()
	w := performRequest(router, "GET", "/", nil)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["hello"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["hello"], value)
}
