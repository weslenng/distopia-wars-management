package handler

import (
	"context"
	"fmt"
	"slices"
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

	possibleTeams := strings.Split(h.cfg.Discord.PossibleTeams, ",")
	playerTeam := strings.ToUpper(args[1])

	if !slices.Contains(possibleTeams, playerTeam) {
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

	if _, err := h.function.Join(ctx, message.Author.ID, playerTeam); err != nil {
		reply := &discordgo.MessageSend{
			Content:   err.Error(),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("join_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.PositiveReaction); err != nil {
		h.log.Error().Err(err).Msgf("join_handler -> failed to add feedback reaction to message %s", message.ID)
	}
}
