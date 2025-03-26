package models

type Player struct {
	ID                int64  `db:"player_id"`
	DiscordID         string `db:"discord_id"`
	MinecraftNickname string `db:"minecraft_nickname"`
	MinecraftPassword string `db:"minecraft_password"`
	MinecraftTeam     string `db:"minecraft_team"`
}
