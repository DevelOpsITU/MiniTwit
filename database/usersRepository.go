package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"minitwit/models"
	"strconv"
)

// TODO: Return errors if any, and meybe the user
func GormAddUserToDb(user models.RegistrationUser) uint {

	salt := make([]byte, 4)
	io.ReadFull(rand.Reader, salt)

	pwIteration_int, _ := strconv.Atoi("50000")
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
		panic(result.Error)
	} else if result.RowsAffected != 1 {
		panic(fmt.Sprint(result.RowsAffected) + " rows were affected")
	}

	return uint(user_obj.UserId)

}

func GormGetUserFromDb(username string) (models.User, error) {
	var user User
	result := gormDb.
		Where("username like ?", username).
		First(&user)

	if result.Error != nil {
		return models.User{}, errors.New("User not found")
		//panic(result.Error)
	}

	return models.User{
		User_id:  user.UserId,
		Username: user.Username,
		Email:    user.Email,
		Pw_hash:  user.PwHash,
	}, nil

}

func getUserFromDb(userId uint) User {
	var user User

	gormDb.Find(&user, userId)

	return user
}

func GormRemoveUserFromDb(user_id uint) {

	result := gormDb.
		Delete(&User{}, user_id)

	if result.Error != nil {
		panic(result.Error)
	} else if result.RowsAffected != 1 {
		panic(fmt.Sprint(result.RowsAffected) + " rows were affected")
	}

}

func NumberOfUsers() int64 {
	var count int64
	gormDb.Model(&User{}).Count(&count)
	return count
}
