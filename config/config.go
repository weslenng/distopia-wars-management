package config

import "github.com/spf13/viper"

type Discord struct {
	GuildID          string
	ChannelID        string
	DeveloperID      string
	RegisterRoleID   string
	SubscriberRoleID string
	Prefix           string
	Token            string
}

type Service struct {
	Debug bool
	Name  string
}

type Config struct {
	Discord *Discord
	Service *Service
}

func New() *Config {
	viper.AutomaticEnv()

	return &Config{
		Discord: &Discord{
			GuildID:          viper.GetString("DISTOPIA_WARS_DISCORD_GUILD_ID"),
			ChannelID:        viper.GetString("DISTOPIA_WARS_DISCORD_CHANNEL_ID"),
			DeveloperID:      viper.GetString("DISTOPIA_WARS_DISCORD_DEVELOPER_ID"),
			RegisterRoleID:   viper.GetString("DISTOPIA_WARS_DISCORD_REGISTER_ROLE_ID"),
			SubscriberRoleID: viper.GetString("DISTOPIA_WARS_DISCORD_SUBSCRIBER_ROLE_ID"),
			Prefix:           viper.GetString("DISTOPIA_WARS_DISCORD_PREFIX"),
			Token:            viper.GetString("DISTOPIA_WARS_DISCORD_TOKEN"),
		},
		Service: &Service{
			Debug: viper.GetBool("DISTOPIA_WARS_SERVICE_DEBUG"),
			Name:  viper.GetString("DISTOPIA_WARS_SERVICE_NAME"),
		},
	}
}
