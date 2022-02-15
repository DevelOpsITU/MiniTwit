package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"golang.org/x/crypto/pbkdf2"
	"minitwit/src/database"
	"minitwit/src/functions"
	"minitwit/src/models"
	"net/http"
	"strconv"
	"strings"
)

var loginTemplate = gonja.Must(gonja.FromFile("templates/login.html"))

func HandleLogin(w gin.ResponseWriter, r *http.Request, c *gin.Context) {

	//var error = ""
	//GetAllMessages()

	if r.Method == http.MethodPost {
		username := c.PostForm("username")
		pw := c.PostForm("password")
		print(pw)
		print(username)
		user := database.GetUserFromDb(username)
		print(user.Pw_hash)
		s := strings.Split(user.Pw_hash, ":")

		s2 := strings.Split(s[2], "$")

		pwIteration := s2[0]
		pwSalt := s2[1]
		pwHash := s2[2]
		fmt.Println(pwIteration)
		fmt.Println(pwSalt)
		fmt.Println(pwHash)

		//salt := []byte(user.pw_hash[21:37]) // TODO: what is this used for?
		pwIteration_int, _ := strconv.Atoi(pwIteration)

		dk := pbkdf2.Key([]byte(pw), []byte(pwSalt), pwIteration_int, 32, sha256.New)

		fmt.Printf("\nsha256: %x\n", []byte(dk))

		// TODO: same as l.293 (What is thi used for?)
		//fmt.Printf("salt: %x\n", string(salt))
		//fmt.Println("len(salt)", len(salt),
		//	"\nlen(hashed)", len(dk))

		if hex.EncodeToString(dk) != pwHash {
			// Invalid
			print("Invalid password")
			c.Redirect(http.StatusFound, "/login")

			return
		} else {
			// User is authenticated

			var g = models.Session{
				User:     user,
				Message:  true,
				Messages: []string{"You were successfully logged in"},
			}
			data, _ := json.Marshal(g)
			c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusFound, "/")
			return
		}

	}
	cookie, err := functions.GetCookie(c)
	if err != nil {
		out, err := loginTemplate.Execute(gonja.Context{"g": ""})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(out))
	} else {
		out, err := loginTemplate.Execute(gonja.Context{"g": cookie})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(out))
	}
}
