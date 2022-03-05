package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"minitwit/log"
	"minitwit/models"
	"strconv"
)

// TODO: Return errors if any, and meybe the user
func AddUserToDb(user models.RegistrationUser) uint {

	salt := make([]byte, 4)
	io.ReadFull(rand.Reader, salt)

	pwIteration_int, _ := strconv.Atoi("50000") //TODO: Hvorfor fanden parser vi en streng ned til et tal? xD
	dk := pbkdf2.Key([]byte(user.Password1), salt, pwIteration_int, 32, sha256.New)

	pw_hashed := "pbkdf2:sha256:50000$" + string(salt) + "$" + hex.EncodeToString(dk)

	user_obj := User{
		Username: user.Username,
		Email:    user.Email,
		PwHash:   pw_hashed,
	}

	result := gormDb.
		Create(&user_obj)

	if result.Error != nil {
		log.Logger.Error().Msg("Could not create the user")
	} else if result.RowsAffected != 1 {
		log.Logger.Error().Str("rowsAffected", fmt.Sprint(result.RowsAffected)).Msg("More rows affected than expected")
	}

	return user_obj.UserId

}

func GetUserFromDb(username string) (User, error) {
	var user User
	result := gormDb.
		Where("username like ?", username).
		First(&user)

	if result.Error != nil {
		return User{}, errors.New("user not found")
	}

	return user, nil

}

func GetUserFromDbWithId(userId uint) User {
	var user User

	gormDb.Find(&user, userId)

	return user
}

func RemoveUserFromDb(userId uint) {

	result := gormDb.
		Delete(&User{}, userId)

	if result.Error != nil {
		log.Logger.Error().Msg("Could not delete the user")
	} else if result.RowsAffected != 1 {
		log.Logger.Error().Str("rowsAffected", fmt.Sprint(result.RowsAffected)).Msg("More rows affected than expected")
	}
}

func NumberOfUsers() int64 {
	var count int64
	gormDb.Model(&User{}).Count(&count)
	return count
}

func CheckIfUserExists(username string) bool {

	user, _ := GetUserFromDb(username)
	if user.UserId != 0 {
		return true
	}

	return false
}
