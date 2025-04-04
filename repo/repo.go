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
	// Create migrations table if it doesn't exist
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			migration_id INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	if _, err := r.sqlite.ExecContext(ctx, query); err != nil {
		return err
	}

	functions := []func(ctx context.Context) error{
		r.Migrate1,
		r.Migrate2,
	}

	for functionIndex, function := range functions {
		migrationID := functionIndex + 1

		// Check if migration has already been applied
		var count int
		if err := r.sqlite.GetContext(ctx, &count, "SELECT COUNT(*) FROM migrations WHERE migration_id = ?", migrationID); err != nil {
			return err
		}

		if count > 0 {
			r.log.Info().Msgf("Migration %d already applied, skipping", migrationID)
			continue
		}

		r.log.Info().Msgf("Applying migration %d of %d", migrationID, len(functions))

		if err := function(ctx); err != nil {
			return err
		}

		// Record that migration was applied
		if _, err := r.sqlite.ExecContext(ctx, "INSERT INTO migrations (migration_id) VALUES (?)", migrationID); err != nil {
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
