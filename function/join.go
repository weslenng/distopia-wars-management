package function

import (
	"context"

	"github.com/weslenng/cuba-wars-management/misc/errors"
	"github.com/weslenng/cuba-wars-management/models"
)

func (f *Function) Join(ctx context.Context, discordID, team string) (player models.Player, err error) {
	player, err = f.repo.SelectPlayerByDiscordID(ctx, discordID)
	if err != nil {
		return player, err
	}

	if player.ID == 0 {
		return player, errors.PlayerNotFound
	}

	if err := f.repo.UpdatePlayerTeam(ctx, player.ID, team); err != nil {
		return player, err
	}

	player.MinecraftTeam = team

	return player, nil
}
