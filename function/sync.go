package function

import (
	"context"

	"github.com/weslenng/cuba-wars-management/models"
)

func (f *Function) Sync(ctx context.Context) (players []models.Player, err error) {
	players, err = f.repo.GetPlayers(ctx)
	if err != nil {
		return players, err
	}

	return players, nil
}
