package logic

import (
	"errors"
	"fmt"
	"minitwit/database"
	"minitwit/functions"
	"minitwit/log"
	"minitwit/models"
	"strings"
)

func CreateUser(registationUser models.RegistrationUser) error {

	if registationUser.Username == "" {
		return errors.New("You have to enter a username")
	} else if registationUser.Email == "" || !strings.Contains(registationUser.Email, "@") {
		return errors.New("Your have to enter a valid email address")
	} else if registationUser.Password1 == "" {
		return errors.New("You have to enter a password")
	} else if registationUser.Password1 != registationUser.Password2 {
		return errors.New("The two passwords do not match")
	} else {

		userExists := database.CheckIfUserExists(registationUser.Username)
		if userExists { // Connection with database error or mapping of user to object
			return errors.New("The username is already taken")
		} else {
			database.AddUserToDb(registationUser)
			return nil
		}
	}
}

func FollowUser(userId uint, usernameToFollow string) error {

	userToFollow, err := database.GetUserFromDb(usernameToFollow)

	if err != nil {
		log.Logger.Error().Err(err).Str("userId", fmt.Sprint(userId)).Msg("Could not get user")
		return err
	}

	//TODO: Check that the user is not already following the user Issue #47
	err = database.FollowUser(userId, userToFollow.UserId)

	if err != nil {
		log.Logger.Error().Err(err).Str("followerId", fmt.Sprint(userId)).Str("followedId", fmt.Sprint(userToFollow.UserId)).Msg("Could not follow user")
	}

	return nil
}

func FollowUserFromUsername(followerUsername string, usernameToFollow string) error {

	follower, err := database.GetUserFromDb(followerUsername)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("username", followerUsername).Msg("Could not get user")
		return err
	}

	userToFollow, err := database.GetUserFromDb(usernameToFollow)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("username", usernameToFollow).Msg("Could not get user")
		return err
	}

	//TODO: Check that the user is not already following the user Issue #47
	err = database.FollowUser(follower.UserId, userToFollow.UserId)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("follower", followerUsername).Str("followed", usernameToFollow).Msg("Could not follow user")
	}

	return nil
}

func IsFollowing(id uint, name string) (bool, error) {
	followerids, err := database.GetFollowingUsers(id)
	if err != nil {
		return false, err
	}
	user, err := database.GetUserFromDb(name)
	if err != nil {
		return false, err
	}
	return functions.ContainsUint(followerids, user.UserId), nil
}

func UnFollowUserFromUsername(followerUsername string, unfollowUsername string) error {

	follower, err := database.GetUserFromDb(followerUsername)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("username", followerUsername).Msg("Could not get user")
		return err
	}

	userToUnFollow, err := database.GetUserFromDb(unfollowUsername)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("username", unfollowUsername).Msg("Could not get user")
		return err
	}

	//TODO: Check that the user is not already following the user Issue #47
	err = database.UnFollowUser(follower.UserId, userToUnFollow.UserId)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("follower", followerUsername).Str("followed", unfollowUsername).Msg("Could not unfollow user")
	}

	return nil
}

func UnFollowUser(userId uint, usernameToUnFollow string) error {
	userToUnFollow, err := database.GetUserFromDb(usernameToUnFollow)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("username", usernameToUnFollow).Msg("Could not get user")
		return err
	}
	// TODO: check if already following before trying this
	err = database.UnFollowUser(userId, userToUnFollow.UserId)

	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("follower", fmt.Sprint(userId)).Str("followed", usernameToUnFollow).Msg("Could not unfollow user")
	}

	return nil

}

func GetUserFollowerUsernames(username string, limit int) ([]string, error) {
	user, err := database.GetUserFromDb(username)
	if err != nil {
		log.Logger.Error().Err(err).Caller().Str("username", username).Msg("Could not get user")
		return nil, err
	}

	users_int, err := database.GetFollowingUsers(user.UserId)
	if err != nil {
		return nil, err
	}

	var usernames = []string{}
	for _, user_id := range users_int {
		usernames = append(usernames, database.GetUserFromDbWithId(user_id).Username)
	}

	return usernames, nil
}
