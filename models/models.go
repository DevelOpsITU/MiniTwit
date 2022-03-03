package models

/****************************************
*		   DATABASE ENTITIES			*
****************************************/

type Message struct {
	MessageId uint   `gorm:"column:message_id;primaryKey;autoIncrement"`
	AuthorId  uint   `gorm:"column:author_id;type:not null"`
	Text      string `gorm:"column:text;type:string not null"`
	Pubdate   int64  `gorm:"column:pub_date"`
	Flagged   int    `gorm:"column:flagged"`
	Username  string
	Email     string
}

type RegistrationUser struct {
	Username  string
	Email     string
	Password1 string
	Password2 string
}

type User struct {
	User_id  uint
	Username string
	Email    string
	Pw_hash  string
}

type Session struct {
	User     User
	Message  bool
	Messages []string
}

type Request struct {
	Endpoint string
}

type Twit struct {
	GavatarUrl string
	Username   string
	Pub_date   string
	Text       string
}
