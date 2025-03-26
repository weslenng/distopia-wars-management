package function

import (
	"context"

	"github.com/weslenng/cuba-wars-management/misc/errors"
	"github.com/weslenng/cuba-wars-management/models"
)

func (f *Function) Login(ctx context.Context, discordID string) (player models.Player, err error) {
	player, err = f.repo.SelectPlayerByDiscordID(ctx, discordID)
	if err != nil {
		return player, err
	}

	if player.ID == 0 {
		return player, errors.PlayerNotFound
	}

	return player, nil
}
