package models

type Message struct {
	MessageId uint
	AuthorId  uint
	Text      string
	Pubdate   int64
	Flagged   int
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
