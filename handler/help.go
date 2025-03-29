package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/weslenng/cuba-wars-management/misc/consts"
	"github.com/weslenng/cuba-wars-management/misc/util"
)

func (h *Handler) getHelpPrefix() string {
	return h.cfg.Discord.Prefix + "help"
}

func (h *Handler) Help(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	if action != h.getHelpPrefix() {
		return
	}

	if h.cfg.Service.Debug {
		if message.Author.ID != h.cfg.Discord.DeveloperID {
			if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
				h.log.Error().Err(err).Msgf("help_handler -> failed to add feedback reaction to message %s", message.ID)
			}

			return
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Distopia Wars 🗡️",
		Description: "Para se registrar digite `.register`\n\nComandos disponíveis:",
		Color:       15132410,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  ".info",
				Value: "Exibe as informações públicas da sua conta",
			},
			{
				Name:  ".join `team`",
				Value: "Se junta a um time, deixando o anterior caso possua",
			},
			{
				Name:  ".login",
				Value: "Envia uma DM com a sua senha do Distopia Wars",
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
		h.log.Error().Err(err).Msgf("help_handler -> failed to send message to channel %s", message.ChannelID)

		if err := session.MessageReactionAdd(message.ChannelID, message.ID, consts.NegativeReaction); err != nil {
			h.log.Error().Err(err).Msgf("help_handler -> failed to add feedback reaction to message %s", message.ID)
		}

		return
	}
}
