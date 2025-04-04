package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

type Repository struct {
	log    zerolog.Logger
	sqlite *sqlx.DB
}

func New(ctx context.Context, log zerolog.Logger) (repo *Repository, err error) {
	sqlite, err := sqlx.Open("sqlite3", "distopia-wars.sqlite3")
	if err != nil {
		return nil, err
	}

	repo = &Repository{
		log:    log,
		sqlite: sqlite,
	}

	if err := repo.Migrate(ctx); err != nil {
		return repo, err
	}

	return repo, nil
}

func (r *Repository) Migrate(ctx context.Context) (err error) {
	functions := []func(ctx context.Context) error{
		r.Migrate1,
		r.Migrate2,
	}

	for functionIndex, function := range functions {
		r.log.Info().Msgf("Migrating %d of %d", functionIndex+1, len(functions))

		if err := function(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Migrate1(ctx context.Context) (err error) {
	query := `
		CREATE TABLE IF NOT EXISTS player (
			player_id INTEGER PRIMARY KEY AUTOINCREMENT,
			discord_id TEXT UNIQUE NOT NULL,
			minecraft_nickname TEXT UNIQUE NOT NULL,
			minecraft_password TEXT NOT NULL,
			minecraft_team TEXT NOT NULL
		);`

	if _, err := r.sqlite.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Migrate2(ctx context.Context) (err error) {
	query := `ALTER TABLE player ADD COLUMN can_change_minecraft_team BOOLEAN DEFAULT 1;`

	if _, err := r.sqlite.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Close() (err error) {
	return r.sqlite.Close()
}
