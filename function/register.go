package function

import (
	"context"
	"strings"

	"github.com/weslenng/cuba-wars-management/misc/errors"
	"github.com/weslenng/cuba-wars-management/misc/util"
	"github.com/weslenng/cuba-wars-management/models"
)

func (f *Function) Register(ctx context.Context, discordID, nickname string) (player models.Player, err error) {
	player, err = f.repo.SelectPlayerByDiscordID(ctx, discordID)
	if err != nil {
		return player, err
	}

	if player.ID != 0 {
		return player, errors.PlayerAlreadyRegistered
	}

	player = models.Player{
		DiscordID:              discordID,
		MinecraftNickname:      nickname,
		MinecraftPassword:      util.NewPassword(),
		CanChangeMinecraftTeam: true,
	}

	playerID, err := f.repo.InsertPlayer(ctx, player)
	if err != nil {
		message := err.Error()

		if strings.Contains(message, "UNIQUE") {
			return player, errors.NicknameAlreadyExists
		}

		return player, err
	}

	player.ID = playerID

	return player, nil
}
