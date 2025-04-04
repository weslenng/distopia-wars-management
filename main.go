package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/weslenng/cuba-wars-management/config"
	"github.com/weslenng/cuba-wars-management/handler"
	"github.com/weslenng/cuba-wars-management/misc/logger"
	"github.com/weslenng/cuba-wars-management/repo"
)

func main() {
	ctx := context.Background()

	cfg := config.New()
	log := logger.New(cfg.Service.Name)

	repo, err := repo.New(ctx, log)
	end(log, err)

	handler, err := handler.New(cfg, log, repo)
	end(log, err)

	discord, err := discordgo.New(cfg.Discord.Token)
	end(log, err)

	discord.AddHandler(handler.Force)
	discord.AddHandler(handler.Help)
	discord.AddHandler(handler.Info)
	discord.AddHandler(handler.Join)
	discord.AddHandler(handler.Login)
	discord.AddHandler(handler.Register)
	discord.AddHandler(handler.Sync)

	discord.Identify.Intents = discordgo.IntentsAll

	if err := discord.Open(); err != nil {
		end(log, err)
	}

	log.Info().Msg("Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := discord.Close(); err != nil {
		end(log, err)
	}
}

func end(log zerolog.Logger, err error) {
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
