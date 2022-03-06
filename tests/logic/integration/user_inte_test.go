package tests

import (
	"bytes"
	"encoding/json"
	"minitwit/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*************************************************
* Register
**************************************************/
func Test_Register_Invalid_Model(t *testing.T) {

	prepare()
	router := setup()
	w := performRequest(router, "POST", "/register", nil)
	cookies := w.Result().Cookies()

	assert.Equal(t, "You have to enter a username", cookies)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_Register_Valid_Model(t *testing.T) {
	data, _ := json.Marshal(models.RegistrationUser{Username: "TEST", Email: "@", Password1: "t", Password2: "t"})
	prepare()
	router := setup()
	w := performRequest(router, "POST", "/register", bytes.NewBuffer(data))

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, http.Cookie{}, w.Result().Cookies()[0])
}
