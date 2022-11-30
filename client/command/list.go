package command

import "github.com/bwmarrin/discordgo"

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "list",
		Description: "View room info and/or delete the room.",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "room",
				Description:  "Select an room to view.",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
			{
				Name:        "delete-room",
				Description: "Delete select room. Only room creator can delete their rooms.",
				Type:        discordgo.ApplicationCommandOptionBoolean,
			},
		},
	})
}
