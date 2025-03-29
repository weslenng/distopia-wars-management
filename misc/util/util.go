package util

import (
	"math/rand"
	"regexp"
	"strings"
)

const Space = " "

func Break(message string) []string {
	message = strings.TrimSpace(message)

	return strings.Fields(message)
}

func Nickname(value string) (nickname string, match bool) {
	nickname = strings.ToLower(value)
	match, _ = regexp.MatchString(`^[a-zA-Z0-9_]{3,16}$`, nickname)

	return nickname, match
}

func NewPassword() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)

	for resultIndex := range result {
		result[resultIndex] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
