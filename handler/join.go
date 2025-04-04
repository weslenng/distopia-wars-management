package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getJoinPrefix() string {
	return h.cfg.Discord.Prefix + "join"
}

const JoinArgs = 2

func (h *Handler) Join(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	if action != h.getJoinPrefix() {
		return
	}

	if h.cfg.Service.Debug {
		if message.Author.ID != h.cfg.Discord.DeveloperID {
			if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
				h.log.Error().Err(err).Msgf("join_handler -> failed to add feedback reaction to message %s", message.ID)
			}

			return
		}
	}

	reference := &discordgo.MessageReference{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
	}

	if len(args) != JoinArgs {
		reply := &discordgo.MessageSend{
			Content:   fmt.Sprintf("O comando deve conter dois argumentos. e.g. .join <@%s> ISRAEL", message.Author.ID),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	playerTeam := strings.ToUpper(args[1])
	possibleTeams := util.ParseTeams(h.cfg.Discord.PossibleTeams)

	possibleTeam := util.Team{}

	for _, entry := range possibleTeams {
		if playerTeam == entry.Name {
			possibleTeam = entry
			break
		}
	}

	if len(possibleTeam.RoleID) == 0 {
		reply := &discordgo.MessageSend{
			Content:   "O time escolhido é inválido",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	ctx := context.Background()

	if _, err := h.function.Join(ctx, message.Author.ID, possibleTeam.Name); err != nil {
		reply := &discordgo.MessageSend{
			Content:   err.Error(),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	for _, possibleTeam := range possibleTeams {
		if err := session.GuildMemberRoleRemove(h.cfg.Discord.GuildID, message.Author.ID, possibleTeam.RoleID); err != nil {
			h.log.Error().Err(err).Msgf("join_handler -> failed to remove role %s from player %s", possibleTeam.RoleID, message.Author.ID)
		}
	}

	if err := session.GuildMemberRoleAdd(h.cfg.Discord.GuildID, message.Author.ID, possibleTeam.RoleID); err != nil {
		h.log.Error().Err(err).Msgf("join_handler -> failed to add role %s to player %s", possibleTeam.RoleID, message.Author.ID)
	}

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.PositiveReaction); err != nil {
		h.log.Error().Err(err).Msgf("join_handler -> failed to add feedback reaction to message %s", message.ID)
	}
}
