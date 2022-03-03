package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"minitwit/database"
	"minitwit/models"
	"strconv"
	"strings"
	"time"
)

func GetUserTwits(username string) ([]models.Twit, models.User, error) {
	user, err := database.GormGetUserFromDb(username)

	if err != nil {
		return []models.Twit{}, models.User{}, err
	} else {
		messages, err := database.GormGetUserMessages(user.User_id)
		if err != nil {
			return []models.Twit{}, models.User{}, err
		}

		return ConvertMessagesToTwits(&messages), user, nil
	}

}

func GetPublicTimelineTwits() ([]models.Twit, error) {
	messages := database.GormGetAllMessages()
	return ConvertMessagesToTwits(&messages), nil

}

func GetPersonalTimelineTwits(user models.User) ([]models.Twit, error) {
	user, err := database.GormGetUserFromDb(user.Username)

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
		twits = append(twits, models.Twit{GavatarUrl: getGavaterUrl(message.Email, 48), Username: message.Username, Pub_date: (formatPubdate(message.Pubdate)), Text: message.Text})
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

func formatPubdate(Pubdate int64) string {
	date := time.Unix(Pubdate, 0)

	formatted := fmt.Sprintf("%d-%02d-%02d @ %02d:%02d",
		date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute())

	return formatted
}
