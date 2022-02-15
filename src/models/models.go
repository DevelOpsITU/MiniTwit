package models

/****************************************
*		   DATABASE ENTITIES			*
****************************************/
type Message struct {
	MessageId int
	AuthorId  int
	Username  string
	Text      string
	Pubdate   int64
	Flagged   bool
	Email     string
}

type User struct {
	User_id  int
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
