package handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getInfoPrefix() string {
	return h.cfg.Discord.Prefix + "info"
}

func (h *Handler) Info(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	if action != h.getInfoPrefix() {
		return
	}

	if h.cfg.Service.Debug {
		if message.Author.ID != h.cfg.Discord.DeveloperID {
			if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
				h.log.Error().Err(err).Msgf("info_handler -> failed to add feedback reaction to message %s", message.ID)
			}

			return
		}
	}

	reference := &discordgo.MessageReference{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
	}

	ctx := context.Background()

	player, err := h.function.Info(ctx, message.Author.ID)
	if err != nil {
		reply := &discordgo.MessageSend{
			Content:   err.Error(),
			Reference: reference,
		}

		if _, err := session.ChannelMessageSendComplex(message.ChannelID, reply); err != nil {
			h.log.Error().Err(err).Msgf("info_handler -> failed to reply to message %s", message.ID)
		}

		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Distopia Wars 🗡️",
		Description: "Informações do jogador",
		Color:       15132410,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Nickname",
				Value: player.MinecraftNickname,
			},
			{
				Name:  "Time",
				Value: player.MinecraftTeam,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: message.Author.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: consts.Footer,
		},
	}

	if _, err := session.ChannelMessageSendEmbed(message.ChannelID, embed); err != nil {
		h.log.Error().Err(err).Msgf("info_handler -> failed to send message to channel %s", message.ChannelID)

		if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
			h.log.Error().Err(err).Msgf("info_handler -> failed to add feedback reaction to message %s", message.ID)
		}

		return
	}
}
