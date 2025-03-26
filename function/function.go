package function

import (
	"github.com/rs/zerolog"
	"github.com/weslenng/cuba-wars-management/config"
	"github.com/weslenng/cuba-wars-management/repo"
)

type Function struct {
	cfg  *config.Config
	log  zerolog.Logger
	repo *repo.Repository
}

func New(cfg *config.Config, log zerolog.Logger, repo *repo.Repository) (function *Function, err error) {
	function = &Function{
		cfg:  cfg,
		log:  log,
		repo: repo,
	}

	return function, nil
}
