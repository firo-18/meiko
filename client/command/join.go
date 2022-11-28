package command

import "github.com/bwmarrin/discordgo"

var (
	MinISV float64 = 60
	MaxISV float64 = 150
)

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "join",
		Description: "Join the current/upcoming event as a runner/filler.",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "room",
				Description:  "Select a room to join.",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
			{
				Name:        "lead",
				Description: "Enter your leader's score up %.",
				Type:        discordgo.ApplicationCommandOptionInteger,
				MinValue:    &MinISV,
				MaxValue:    MaxISV,
				Required:    true,
			},
			{
				Name:        "sum",
				Description: "Enter the sum ISV total.",
				Type:        discordgo.ApplicationCommandOptionInteger,
				MinValue:    &MinISV,
				MaxValue:    MaxISV * 5,
				Required:    true,
			},
		},
	})
}
