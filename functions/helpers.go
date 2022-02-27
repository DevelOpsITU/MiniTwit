package functions

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minitwit/models"
	config "minitwit/settings"
	"net/http"
)

func GetCookie(c *gin.Context) (models.Session, error) {
	var g models.Session
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
	SetCookie(c, newCookie)

	return g, nil

}

func SetCookie(c *gin.Context, session models.Session) {

	data, _ := json.Marshal(session)
	c.SetCookie("session", string(data), 3600, "/", config.GetConfig().Server.Host, false, true)
}

func GetEndpoint(r *http.Request) models.Request {
	var request = models.Request{r.URL.Path}

	if request.Endpoint == "/public" {
		request.Endpoint = "public_timeline"
	} else if len(request.Endpoint) > 1 {
		request.Endpoint = "user_timeline"
	} else {
		request.Endpoint = "timeline"
	}
	return request
}
