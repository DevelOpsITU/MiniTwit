package logic

import (
	"testing"
)

func TestGavatarString(t *testing.T) {
	email := "user@mail.com"
	size := "48"
	sizeInt := 48
	hashStr := "6ad193f57f79ac444c3621370da955e9"
	gavatarStr := getGavaterUrl(email, sizeInt)

	expectedString := "http://www.gravatar.com/avatar/" + hashStr + "?d=identicon&s=" + size

	if gavatarStr != expectedString {
		t.Fatalf("Gavtar failed")
	}
}

// https://docs.coveralls.io/go
