package logic

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"minitwit/database"
	"minitwit/models"
	"strconv"
	"strings"
)

func CheckPassword(username string, pw string) (models.User, error) {
	user, err := database.GetUserFromDb(username)

	if err != nil {
		return models.User{}, errors.New("invalid username")
	}

	s := strings.Split(user.PwHash, ":")

	s2 := strings.Split(s[2], "$")

	pwIteration := s2[0]
	pwSalt := s2[1]
	pwHash := s2[2]

	passwordIterationInt, _ := strconv.Atoi(pwIteration)

	dk := pbkdf2.Key([]byte(pw), []byte(pwSalt), passwordIterationInt, 32, sha256.New)

	if hex.EncodeToString(dk) == pwHash {
		return models.User{
			User_id:  user.UserId,
			Username: user.Username,
		}, nil
	} else {
		return models.User{}, errors.New("invalid password")
	}
}
