package database

import (
	"errors"
)

func contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func FollowUser(userId uint, UserIdToFollow uint) error {

	follower := GetUserFromDbWithId(userId)
	followed := GetUserFromDbWithId(UserIdToFollow)

	if follower.UserId == 0 || followed.UserId == 0 {
		return errors.New("user did not exist")
	}

	if contains(GetFollowingUsers(follower.UserId), followed.UserId) {
		return errors.New("user already follows given user")
	}

	obj := Follower{
		WhoId:  userId,
		WhomId: UserIdToFollow,
	}

	create := gormDb.Model(&Follower{}).Create(obj)

	if create.Error != nil {
		println(create.Error.Error())
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
		return errors.New("error when unfollowing user")
	}
	return nil
}

// Returns a list of all the users a user is following
func GetFollowingUsers(userId uint) []uint {

	var follows []uint

	subquery, err := gormDb.
		Model(&Follower{}).
		Select("whom_id").
		Where("who_id = ?", userId).
		Rows()

	if err != nil {
		//TODO: Remove panic statements. it crashes the application.
		panic(err)
	}

	for subquery.Next() {
		var user uint
		err := subquery.Scan(&user)
		if err != nil {
			//TODO
		}
		follows = append(follows, user)
	}
	return follows
}
