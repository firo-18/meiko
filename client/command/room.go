package command

import "github.com/bwmarrin/discordgo"

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "room",
		Description: "Create a new tiering room as a runner.",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "name",
				Description: "Enter a room name. Meme is fine, as long as people know which room is for who.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:         "event",
				Description:  "Select an event to tier.",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
		},
	})
}
