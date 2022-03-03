package database

import (
	"errors"
)

func FollowUser(userId uint, UserIdToFollow uint) error {

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
