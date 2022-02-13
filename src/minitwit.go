package main

import (
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/noirbizarre/gonja"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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
	pw_hash  string
}

const DATABASE = "/tmp/minitwit.db"

//const DATABASE = "C:/Users/hardk/source/repos/MiniTwit/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

/****************************************
*			DATABASE RELATED			*
****************************************/
func ConnectDb() *sql.DB {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}

	return db
}

// setup
func InitDb() {
	db := ConnectDb()
	query, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

// example Database usage
func GetAllMessages() []Message {
	db := ConnectDb()
	query := string("select message.message_id , message.author_id , user.username , message.text , message.pub_date ,  user.email from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit 30")
	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var messages []Message

	for result.Next() {
		var msg Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			panic(err.Error())
		}
		messages = append(messages, msg)
	}
	return messages
}

func GetUserFromDb(username string) User {
	db := ConnectDb()
	//TODO: Prepared statements
	strs := []string{"SELECT x.* FROM 'user' x WHERE username like '", username, "'"}
	query := strings.Join(strs, "")
	row, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var user User
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&user.User_id, &user.Username, &user.Email, &user.pw_hash)
	}

	return user

}

/****************************************
*			REST OF PROGRAM				*
****************************************/

func getCookie(c *gin.Context) (Session, error) {
	var g Session
	cookie, err := c.Cookie("session")

	// If there is no cookie
	if err != nil {
		return g, err
	} else {
		//data,_ := json.Marshal(g)
		//c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
		json.Unmarshal([]byte(cookie), &g)
		//println("Found Cookie:", string([]byte(cookie)))

	}

	newCookie := g
	newCookie.Message = false
	newCookie.Messages = nil
	setCookie(c, newCookie)

	return g, nil

}

func setCookie(c *gin.Context, session Session) {

	data, _ := json.Marshal(session)
	c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
}

// Pre-compiling the templates at application startup using the
// little Must()-helper function (Must() will panic if FromFile()
// or FromString() will return with an error - that's it).
// It's faster to pre-compile it anywhere at startup and only
// execute the template later.

var timelineTemplate = gonja.Must(gonja.FromFile("templates/timeline.html"))
var loginTemplate = gonja.Must(gonja.FromFile("templates/login.html"))

// Route /
func handleTimeline(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// Execute the template per HTTP request
	request := getEndpoint(r)
	var g Session
	g, err := getCookie(c)

	// If there is no cookie
	if err != nil || g.User.Username == "" {
		c.Redirect(http.StatusFound, "/public")
	}

	//set g = "None" if g.user should return false in jinja

	out, err := timelineTemplate.Execute(gonja.Context{"g": g, "request": request})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func getEndpoint(r *http.Request) Request {
	var request = Request{r.URL.Path}

	if request.Endpoint == "/public" {
		request.Endpoint = "public_timeline"
	} else if len(request.Endpoint) > 1 {
		request.Endpoint = "user_timeline"
	} else {
		request.Endpoint = ""
	}
	return request
}

func handlePublicTimeline(w gin.ResponseWriter, r *http.Request, c *gin.Context) {
	request := getEndpoint(r)

	messages := GetAllMessages()

	twits := CovertMessagesToTwits(&messages)
	//print(string(request))
	out, err := timelineTemplate.Execute(gonja.Context{"g": "", "request": request, "messages": twits})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func CovertMessagesToTwits(messages *[]Message) []Twit {
	var twits []Twit
	for _, message := range *messages {
		twits = append(twits, Twit{getGavaterUrl(message.Email, 48), message.Username, strconv.Itoa(int(message.Pubdate)), message.Text})
	}
	print(twits)
	return twits

}

func getGavaterUrl(email string, size int) string {
	data := []byte(strings.ToLower(strings.TrimSpace(email)))
	hash := md5.Sum(data)
	hashStr := hex.EncodeToString(hash[:])

	str := []string{"http://www.gravatar.com/avatar/", hashStr, "?d=identicon&s=", strconv.Itoa(size)}
	return strings.Join(str, "")
}

func handleLogin(w gin.ResponseWriter, r *http.Request, c *gin.Context) {

	//var error = ""
	//GetAllMessages()

	if r.Method == http.MethodPost {
		username := c.PostForm("username")
		pw := c.PostForm("password")
		print(pw)
		print(username)
		user := GetUserFromDb(username)
		print(user.pw_hash)

		s := strings.Split(user.pw_hash, ":")

		s2 := strings.Split(s[2], "$")

		pwIteration := s2[0]
		pwSalt := s2[1]
		pwHash := s2[2]
		fmt.Println(pwIteration)
		fmt.Println(pwSalt)
		fmt.Println(pwHash)

		salt := []byte(user.pw_hash[21:37])
		pwIteration_int, _ := strconv.Atoi(pwIteration)

		dk := pbkdf2.Key([]byte(pw), []byte(pwSalt), pwIteration_int, 32, sha256.New)

		fmt.Printf("\nsha256: %x\n", []byte(dk))
		fmt.Printf("salt: %x\n", string(salt))
		fmt.Println("len(salt)", len(salt),
			"\nlen(hashed)", len(dk))

		if hex.EncodeToString(dk) != pwHash {
			// Invalid
			print("Invalid password")
			c.Redirect(http.StatusFound, "/login")

			return
		} else {
			// User is authenticated

			var g = Session{
				User:     user,
				Message:  true,
				Messages: []string{"You were successfully registered and can login now"},
			}
			data, _ := json.Marshal(g)
			c.SetCookie("session", string(data), 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusFound, "/")
			return
		}

	}

	out, err := loginTemplate.Execute(gonja.Context{"g": ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pongA",
		})
	})

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		handleTimeline(c.Writer, c.Request, c)
	})
	router.GET("/public", func(c *gin.Context) {
		handlePublicTimeline(c.Writer, c.Request, c)
	})

	router.GET("/login", func(c *gin.Context) {
		handleLogin(c.Writer, c.Request, c)
	})

	router.POST("/login", func(c *gin.Context) {
		handleLogin(c.Writer, c.Request, c)
	})

	router.LoadHTMLFiles("./src/test.html")

	/*
	 FOR TESTING GO TOOL 'FRESH': 'go install github.com/pilu/fresh'
	 TRY TO RUN COMMAND: 'fresh -c my_fresh_runner.conf' AND
	 THEN MAKE CHANGES TO THE 'test.html' OR 'minitwit.go' FILES.
	 IF NO ERROR, THEN FRESH SHOULD BUILD AND RUN THE 'minitwit.go' CODE.
	 THE CHANGES SHOULD BE SEEN REFLECTED ON 'http://localhost:8080/test/test.html'.

	 OBS: MAYBE TURN OFF AUTO-SAVING, SO STUFF IS ONLY BUILD AND RAN, WHEN YOU WANT IT TO.
	*/
	router.Static("/test", "./src")

	router.Run(":8080")
	//router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
