package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getLoginPrefix() string {
	return h.cfg.Discord.Prefix + "login"
}

func (h *Handler) Login(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	if action != h.getLoginPrefix() {
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

	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		reply := &discordgo.MessageSend{
			Content:   "Não consegui enviar uma mensagem para você",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("login_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	if _, err := session.ChannelMessageSend(channel.ID, "`...`"); err != nil {
		reply := &discordgo.MessageSend{
			Content:   "Não consegui enviar uma mensagem para você",
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("login_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	ctx := context.Background()

	player, err := h.function.Login(ctx, message.Author.ID)
	if err != nil {
		reply := &discordgo.MessageSend{
			Content:   err.Error(),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("login_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Distopia Wars 🗡️",
		Description: "Informações da conta",
		Color:       15132410,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Senha",
				Value: fmt.Sprintf("`/login %s`", player.MinecraftPassword),
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: message.Author.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: consts.Footer,
		},
	}

	if _, err := session.ChannelMessageSendEmbed(channel.ID, embed); err != nil {
		h.log.Error().Err(err).Msgf("login_handler -> failed to send direct message to user %s", message.Author.ID)

		if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
			h.log.Error().Err(err).Msgf("login_handler -> failed to add feedback reaction to message %s", message.ID)
		}

		return
	}

	if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.PositiveReaction); err != nil {
		h.log.Error().Err(err).Msgf("login_handler -> failed to add feedback reaction to message %s", message.ID)
	}
}
