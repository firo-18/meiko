package command

import (
	"github.com/bwmarrin/discordgo"
)

var (
	MinISV float64 = 60
	MaxISV float64 = 150
)

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "link",
		Description: "Link your account, ISV, and UTC offset info. Re-run this to update info.",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
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
			{
				Name:         "utc-offset",
				Description:  "Enter your UTC offset for local time conversion.",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
		},
	})
}
