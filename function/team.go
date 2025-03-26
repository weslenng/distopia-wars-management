package function

import (
	"context"

	"github.com/weslenng/cuba-wars-management/models"
)

func (f *Function) Team(ctx context.Context, team string) (players []models.Player, err error) {
	players, err = f.repo.GetPlayersByTeam(ctx, team)
	if err != nil {
		return players, err
	}

	return players, nil
}
