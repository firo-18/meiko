package command

import "github.com/bwmarrin/discordgo"

var (
	MinShift float64 = 1
	MaxShift float64 = 39
)

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "schedule",
		Description: "Schedule your availability for a room in 1 hour block.",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "room",
				Description:  "Select a room to schedule.",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
		},
	})
}
