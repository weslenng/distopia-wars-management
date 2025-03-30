package handler

import (
	"context"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getSyncPrefix() string {
	return h.cfg.Discord.Prefix + "sync"
}

func (h *Handler) Sync(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	if message.GuildID != h.cfg.Discord.GuildID {
		return
	}

	if message.ChannelID != h.cfg.Discord.ChannelID {
		return
	}

	args := util.Break(message.Content)

	if len(args) <= 0 {
		return
	}

	action := args[0]

	if action != h.getSyncPrefix() {
		return
	}

	if h.cfg.Service.Debug {
		if message.Author.ID != h.cfg.Discord.DeveloperID {
			if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
				h.log.Error().Err(err).Msgf("sync_handler -> failed to add feedback reaction to message %s", message.ID)
			}

			return
		}
	}

	reference := &discordgo.MessageReference{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
	}

	operator := slices.Contains(message.Member.Roles, h.cfg.Discord.OperatorRoleID)

	if !operator {
		reply := &discordgo.MessageSend{
			Content:   "Você não pode usar este comando.",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("sync_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.SyncReaction); err != nil {
		h.log.Error().Err(err).Msgf("sync_handler -> failed to add feedback reaction to message %s", message.ID)
	}

	ctx := context.Background()

	players, err := h.function.Sync(ctx)
	if err != nil {
		h.log.Error().Err(err).Msgf("sync_handler -> failed to sync players")
	}

	roleTeams := util.ParseTeams(h.cfg.Discord.PossibleTeams)

	if len(roleTeams) != 2 {
		reply := &discordgo.MessageSend{
			Content:   "O número de times é inválido",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("sync_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	for _, player := range players {
		if player.MinecraftTeam == roleTeams[0].Name {
			session.GuildMemberRoleAdd(message.GuildID, player.DiscordID, roleTeams[0].RoleID)
		}

		if player.MinecraftTeam == roleTeams[1].Name {
			session.GuildMemberRoleAdd(message.GuildID, player.DiscordID, roleTeams[1].RoleID)
		}

		time.Sleep(1 * time.Second)
	}

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.PositiveReaction); err != nil {
		h.log.Error().Err(err).Msgf("sync_handler -> failed to add feedback reaction to message %s", message.ID)
	}
}
