package component

import "github.com/bwmarrin/discordgo"

var (
	List = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)
