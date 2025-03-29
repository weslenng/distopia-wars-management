package handler

import (
	"context"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getRegisterPrefix() string {
	return h.cfg.Discord.Prefix + "register"
}

const RegisterArgs = 2

func (h *Handler) Register(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	if action != h.getRegisterPrefix() {
		return
	}

	if h.cfg.Service.Debug {
		if message.Author.ID != h.cfg.Discord.DeveloperID {
			if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
				h.log.Error().Err(err).Msgf("login_handler -> failed to add feedback reaction to message %s", message.ID)
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
	subscriber := slices.Contains(message.Member.Roles, h.cfg.Discord.SubscriberRoleID)

	if !operator {
		if !subscriber {
			reply := &discordgo.MessageSend{
				Content:   "Você não é inscrito no canal. Talvez tenha esquecido de vincular a conta da Twitch com o Discord?",
				Reference: reference,
			}

			if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
				h.log.Error().Err(err).Msgf("register_handler -> failed to reply to message %s", message.ID)
			}

			return
		}
	}

	if len(args) != RegisterArgs {
		reply := &discordgo.MessageSend{
			Content:   "O comando deve conter um argumento. e.g. .register savi2w",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("register_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	nickname, match := util.Nickname(args[1])

	if !match {
		reply := &discordgo.MessageSend{
			Content:   "O nickname escolhido é inválido. Ele deve conter entre 3 e 16 caracteres e pode conter apenas letras, números e sublinhados",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("register_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	ctx := context.Background()

	if _, err := h.function.Register(ctx, message.Author.ID, nickname); err != nil {
		reply := &discordgo.MessageSend{
			Content:   err.Error(),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("register_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	// if err := session.GuildMemberRoleAdd(h.cfg.Discord.GuildID, message.Author.ID, h.cfg.Discord.RegisterRoleID); err != nil {
	// 	h.log.Error().Err(err).Msgf("register_handler -> failed to add role %s to player %s", h.cfg.Discord.RegisterRoleID, message.Author.ID)
	// }

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.PositiveReaction); err != nil {
		h.log.Error().Err(err).Msgf("register_handler -> failed to add feedback reaction to message %s", message.ID)
	}
}
