package handler

import (
	"github.com/rs/zerolog"
	"github.com/weslenng/cuba-wars-management/config"
	"github.com/weslenng/cuba-wars-management/function"
	"github.com/weslenng/cuba-wars-management/repo"
)

type Handler struct {
	cfg      *config.Config
	log      zerolog.Logger
	repo     *repo.Repository
	function *function.Function
}

func New(cfg *config.Config, log zerolog.Logger, repo *repo.Repository) (handler *Handler, err error) {
	function, err := function.New(cfg, log, repo)
	if err != nil {
		return nil, err
	}

	handler = &Handler{
		cfg:      cfg,
		log:      log,
		repo:     repo,
		function: function,
	}

	return handler, nil
}
