package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getForcePrefix() string {
	return h.cfg.Discord.Prefix + "force"
}

const ForceArgs = 3

func (h *Handler) Force(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	if action != h.getForcePrefix() {
		return
	}

	if h.cfg.Service.Debug {
		if message.Author.ID != h.cfg.Discord.DeveloperID {
			if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
				h.log.Error().Err(err).Msgf("force_handler -> failed to add feedback reaction to message %s", message.ID)
			}

			return
		}
	}

	reference := &discordgo.MessageReference{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
	}

	allow := false

	for _, role := range message.Member.Roles {
		if role == h.cfg.Discord.OperatorRoleID {
			allow = true
			break
		}
	}

	if !allow {
		reply := &discordgo.MessageSend{
			Content:   "Você não tem permissão para usar este comando.",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("force_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	if len(args) != ForceArgs {
		reply := &discordgo.MessageSend{
			Content:   fmt.Sprintf("O comando deve conter dois argumentos. e.g. .force <@%s> <team>", message.Author.ID),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("force_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	discordID, match := util.Snowflake(args[1])

	if !match {
		reply := &discordgo.MessageSend{
			Content:   "O ID do usuário é inválido.",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("force_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	playerTeam := strings.ToUpper(args[2])
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
			h.log.Error().Err(err).Msgf("force_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	ctx := context.Background()

	if _, err := h.function.Force(ctx, discordID, possibleTeam.Name); err != nil {
		reply := &discordgo.MessageSend{
			Content:   err.Error(),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("force_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	for _, possibleTeam := range possibleTeams {
		if err := session.GuildMemberRoleRemove(h.cfg.Discord.GuildID, discordID, possibleTeam.RoleID); err != nil {
			h.log.Error().Err(err).Msgf("force_handler -> failed to remove role %s from player %s", possibleTeam.RoleID, discordID)
		}
	}

	if err := session.GuildMemberRoleAdd(h.cfg.Discord.GuildID, discordID, possibleTeam.RoleID); err != nil {
		h.log.Error().Err(err).Msgf("force_handler -> failed to add role %s to player %s", possibleTeam.RoleID, discordID)
	}

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.PositiveReaction); err != nil {
		h.log.Error().Err(err).Msgf("force_handler -> failed to add feedback reaction to message %s", message.ID)
	}
}
