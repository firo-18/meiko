package command

import "github.com/bwmarrin/discordgo"

func init() {
	Private = append(Private, &discordgo.ApplicationCommand{
		Name:                     "stats",
		Description:              "View bot stats.",
		Type:                     discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: &permissionManageServer,
	})
}
