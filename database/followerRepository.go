package database

import (
	"errors"
)

func FollowUser(userId uint, UserIdToFollow uint) error {

	follower := getUserFromDb(userId)
	followed := getUserFromDb(UserIdToFollow)

	if follower.UserId == 0 || followed.UserId == 0 {
		return errors.New("user did not exist")
	}

	obj := Follower{
		WhoId:  uint(userId),
		WhomId: uint(UserIdToFollow),
	}

	create := gormDb.Model(&Follower{}).Create(obj)

	if create.Error != nil {
		println(create.Error.Error())
		return errors.New(create.Error.Error())
	}

	return nil
}
