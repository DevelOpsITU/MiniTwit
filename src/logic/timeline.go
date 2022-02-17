package logic

import (
	"crypto/md5"
	"encoding/hex"
	"minitwit/src/database"
	"minitwit/src/models"
	"strconv"
	"strings"
)

func GetUserTwits(username string) ([]models.Twit, models.User, error) {
	user, err := database.GetUserFromDb(username)

	if err != nil {
		return []models.Twit{}, models.User{}, err
	} else {
		messages := database.GetUserMessages(user.User_id)
		return ConvertMessagesToTwits(&messages), user, nil
	}

}

func GetPersonalTimelineTwits(user models.User) ([]models.Twit, error) {
	user, err := database.GetUserFromDb(user.Username)

	if err != nil {
		return []models.Twit{}, err
	} else {
		messages := database.GetPersonalTimelineMessages(user.User_id)
		return ConvertMessagesToTwits(&messages), nil
	}

}

func ConvertMessagesToTwits(messages *[]models.Message) []models.Twit {
	var twits []models.Twit
	for _, message := range *messages {
		twits = append(twits, models.Twit{getGavaterUrl(message.Email, 48), message.Username, strconv.Itoa(int(message.Pubdate)), message.Text})
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
