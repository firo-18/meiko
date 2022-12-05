package command

import "github.com/bwmarrin/discordgo"

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "manage",
		Description: "Manage a room. You can delete, set a manager, and general management.",
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
				Name:        "manager",
				Description: "Select a user to appoints them as the room manager. They will full access to the room.",
				Type:        discordgo.ApplicationCommandOptionUser,
			},
			{
				Name:        "delete-room",
				Description: "Delete select room. Only room creator can delete their rooms.",
				Type:        discordgo.ApplicationCommandOptionBoolean,
			},
		},
	})
}
