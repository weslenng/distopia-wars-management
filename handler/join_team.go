package handler

import (
	"context"
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getJoinTeamPrefix() string {
	return h.cfg.Discord.Prefix + "join_team"
}

const JoinTeamArgs = 3

func (h *Handler) JoinTeam(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	if action != h.getJoinTeamPrefix() {
		return
	}

	if h.cfg.Service.Debug {
		if message.Author.ID != h.cfg.Discord.DeveloperID {
			if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
				h.log.Error().Err(err).Msgf("join_team_handler -> failed to add feedback reaction to message %s", message.ID)
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
			Content:   "Você não tem permissão para usar este comando",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_team_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	if len(args) != JoinTeamArgs {
		reply := &discordgo.MessageSend{
			Content:   fmt.Sprintf("O comando deve conter dois argumentos. e.g. .join_team <@%s> ISRAEL", message.Author.ID),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_team_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	discordID, match := util.Snowflake(args[1])

	if !match {
		reply := &discordgo.MessageSend{
			Content:   "O ID do usuário é inválido",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_team_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	team, match := util.Team(args[2])

	if !match {
		reply := &discordgo.MessageSend{
			Content:   "O time escolhido é inválido",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_team_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	ctx := context.Background()

	if _, err := h.function.JoinTeam(ctx, discordID, team); err != nil {
		reply := &discordgo.MessageSend{
			Content:   err.Error(),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_team_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.PositiveReaction); err != nil {
		h.log.Error().Err(err).Msgf("join_team_handler -> failed to add feedback reaction to message %s", message.ID)
	}
}
