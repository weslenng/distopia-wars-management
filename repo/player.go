package repo

import (
	"context"

	"github.com/weslenng/cuba-wars-management/models"
)

func (r *Repository) SelectPlayerByDiscordID(ctx context.Context, discordID string) (result models.Player, err error) {
	query := `SELECT * FROM player WHERE discord_id = ?;`

	rows, err := r.sqlite.QueryxContext(ctx, query, discordID)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&result); err != nil {
			return result, err
		}
	}

	return result, nil
}

func (r *Repository) UpdatePlayerTeam(ctx context.Context, playerID int64, team string) (err error) {
	query := `UPDATE player SET minecraft_team = ? WHERE player_id = ?;`

	if _, err := r.sqlite.ExecContext(ctx, query, team, playerID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertPlayer(ctx context.Context, player models.Player) (playerID int64, err error) {
	query := `INSERT INTO player (discord_id, minecraft_nickname, minecraft_password, minecraft_team) VALUES (:discord_id, :minecraft_nickname, :minecraft_password, :minecraft_team);`

	result, err := r.sqlite.NamedExecContext(ctx, query, player)
	if err != nil {
		return playerID, err
	}

	playerID, err = result.LastInsertId()
	if err != nil {
		return playerID, err
	}

	return playerID, nil
}

func (r *Repository) GetPlayersByTeam(ctx context.Context, team string) (players []models.Player, err error) {
	query := `SELECT * FROM player WHERE minecraft_team = ?;`

	rows, err := r.sqlite.QueryxContext(ctx, query, team)
	if err != nil {
		return players, err
	}

	defer rows.Close()

	for rows.Next() {
		var player models.Player

		if err := rows.StructScan(&player); err != nil {
			return players, err
		}

		players = append(players, player)
	}

	return players, nil
}
