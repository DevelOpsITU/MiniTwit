package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"minitwit/config"
)

var gormDb *gorm.DB

type User struct {
	UserId   uint   `gorm:"column:user_id;primaryKey;autoIncrement"`
	Username string `gorm:"column:username;type:string not null"`
	Email    string `gorm:"column:email;type:string not null"`
	PwHash   string `gorm:"column:pw_hash;type:string not null"`
}

func (User) TableName() string {
	return "user"
}

// NOTE: The autoIncrement is not added. https://gorm.io/docs/composite_primary_key.html#content-inner
// Maybe add a post migration script that adds it or live with it?.

type Follower struct {
	WhoId  uint `gorm:"column:who_id"`
	WhomId uint `gorm:"column:whom_id"`
}

func (Follower) TableName() string {
	return "follower"
}

type Message struct {
	MessageId       uint   `gorm:"column:message_id;primaryKey;autoIncrement"`
	AuthorId        uint   `gorm:"column:author_id;type:not null"`
	Text            string `gorm:"column:text;type:string not null"`
	PublicationDate uint   `gorm:"column:pub_date"`
	Flagged         int    `gorm:"column:flagged"`
}

func (Message) TableName() string {
	return "message"
}

func InitGorm() (db *gorm.DB, err error) {

	db, err = gorm.Open(sqlite.Open(config.GetConfig().Database.ConnectionString), &gorm.Config{})
	gormDb = db

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &Follower{}, &Message{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
