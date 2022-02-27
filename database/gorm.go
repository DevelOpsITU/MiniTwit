package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"minitwit/models"
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

	db, err = gorm.Open(sqlite.Open("/tmp/minitwit.db"), &gorm.Config{})
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

func GormGetAllMessages() []models.Message {

	result, err := gormDb.
		Model(models.Message{}).
		Table("message").
		Order("pub_date desc").
		Limit(30).
		Where("flagged = ?", 0).
		Joins("JOIN user on message.author_id = user.user_id").
		Select("message.message_id , message.author_id , user.username , message.text , message.pub_date , user.email").
		Rows()

	if err != nil {
		panic(err)
	}

	var messages2 []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages2 = append(messages2, msg)
	}

	return messages2
}
