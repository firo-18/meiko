package command

import "github.com/bwmarrin/discordgo"

func init() {
	List = append(List, &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Private lesson with Meiko.",
		Type:        discordgo.ChatApplicationCommand,
	})
}
