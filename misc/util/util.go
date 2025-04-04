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

type Team struct {
	RoleID string
	Name   string
}

func ParseTeams(possibleTeams string) (result []Team) {
	roleTeams := strings.Split(possibleTeams, ",")

	for _, roleTeam := range roleTeams {
		args := strings.Split(roleTeam, "=")

		if len(args) != 2 {
			continue
		}

		result = append(result, Team{
			RoleID: args[0],
			Name:   args[1],
		})
	}

	return result
}

func Snowflake(value string) (discordID string, match bool) {
	builder := strings.Builder{}

	for _, char := range value {
		if char >= '0' && char <= '9' {
			builder.WriteRune(char)
		}
	}

	discordID = builder.String()
	match = len(discordID) >= 17 && len(discordID) <= 19

	return discordID, match
}
