package database

import (
	"errors"
	"fmt"
	"minitwit/functions"
	"minitwit/log"
)

func FollowUser(userId uint, UserIdToFollow uint) error {

	follower := GetUserFromDbWithId(userId)
	followed := GetUserFromDbWithId(UserIdToFollow)

	if follower.UserId == 0 || followed.UserId == 0 {
		return errors.New("user did not exist")
	}

	followerids, err := GetFollowingUsers(follower.UserId)
	if err != nil {
		log.Logger.Error().Err(err).Str("userId", fmt.Sprint(follower.UserId)).Msg("Could not get who the user follows")
		return errors.New(err.Error())
	}

	if functions.ContainsUint(followerids, followed.UserId) {
		log.Logger.Error().Err(err).Str("userId", fmt.Sprint(follower.UserId)).Msg("user already follows given user")
		return errors.New("user already follows given user")
	}

	obj := Follower{
		WhoId:  userId,
		WhomId: UserIdToFollow,
	}

	create := gormDb.Model(&Follower{}).Create(obj)

	if create.Error != nil {
		return errors.New(create.Error.Error())
	}

	return nil
}

func UnFollowUser(userId uint, UserIdToUnFollow uint) error {

	obj := Follower{
		WhoId:  userId,
		WhomId: UserIdToUnFollow,
	}

	result := gormDb.
		Delete(&Follower{}, obj)

	if result.RowsAffected != 1 {
		err := errors.New("Could not remove the follower, user is not following the given user")
		log.Logger.Error().Err(err).Str("userId", fmt.Sprint(userId)).Msg(
			fmt.Sprintf("Could not remove follower as user is not following the given user in id: %d", UserIdToUnFollow))
		return err
	}
	return nil
}

// Returns a list of all the users a user is following
func GetFollowingUsers(userId uint) ([]uint, error) {

	var follows []uint

	subquery, err := gormDb.
		Model(&Follower{}).
		Select("whom_id").
		Where("who_id = ?", userId).
		Rows()

	if err != nil {
		return follows, errors.New("error getting following users")
	}

	for subquery.Next() {
		var user uint
		err := subquery.Scan(&user)
		if err != nil {
			return follows, errors.New("error mapping user")
		}
		follows = append(follows, user)
	}
	return follows, nil
}
