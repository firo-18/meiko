package command

import "github.com/bwmarrin/discordgo"

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "session",
		Description: "Begin session to enable ping check in. End the session when done.",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "room",
				Description:  "Select a room.",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
			{
				Name:        "state",
				Description: "Begin or end a session.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Begin",
						Value: "begin",
					},
					{
						Name:  "End",
						Value: "end",
					},
				},
			},
			{
				Name:        "role",
				Description: "Select a role to ping for additional fillers if needed. Can be left empty.",
				Type:        discordgo.ApplicationCommandOptionRole,
			},
		},
	})
}
