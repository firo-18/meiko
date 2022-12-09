package event

import (
	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

func init() {
	List["help"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Welcome to Meiko's Tiering Cafe",
						Description: `Meiko is a tiering management bot. The idea is to have Meiko automates check in process for fillers, sorts fillers by their ISV based on event type, and sends out pings accordingly until event/session ends.
						
						Currently, Meiko has the following functions:`,
						Color:     discord.EmbedColor,
						Timestamp: discord.EmbedTimestamp,
						Footer:    discord.EmbedFooter(s),
						Fields: []*discordgo.MessageEmbedField{
							{
								Name: "/link",
								Value: discord.FieldStyle(`Users need to first link the account first to use other commands, as the infomation inputted in the linking process such as offset and ISV are required to display the correct time format and sorting accordingly.
								
								You can relink as many time you need if your locale or ISV changes.`),
							},
							{
								Name:  "/schedule",
								Value: discord.FieldStyle("Users who has linked can select an active to schedule their shift. This command is straightforward, as you simply click the day and any/all the shifts you want to run/fill. Click outside the select menu will update your shift, if you made any changes."),
							},
							{
								Name: "/room",
								Value: discord.FieldStyle(`A command to create a new room, preferably by the runner. Input a room name (no duplication), and select an event for the room.
								
								Event list is populated by using Sekai database, so if you don't see your event, it means it hasn't been officially announce in EN. I can't do anything about this right now. Fill-all is basically a short cut used to automatically fill all shift with your name (as the runner), and you can manually deselect any shift you don't want later.`),
							},
							{
								Name:  "/view",
								Value: discord.FieldStyle(`View shows the room information as well as it's schedule, by day. If you are in any shift, your name will be underscored.`),
							},
							{
								Name:  "/manage",
								Value: discord.FieldStyle(`Here, the owner/manager of the room can delete the room, appoint a manager, and manage the shift. For shift management, it's only possible to remove fillers from a shift, not assign them, for obvious reason.`),
							},
							{
								Name: "/session",
								Value: discord.FieldStyle(`For sending out check in messages and sorting fillers. Begin the session will send out message hourly (15 mins prior to shift) to ping users who signed up. She will only ping those who just start a new shift.
								
								When you begin the session, you can select a role to ping whenever a shift does not have enough fillers. And don't forget to end the session when you finish the tiering session to avoid pinging unneccessarily.`),
							},
						},
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			ErrExit(err)
		}
	}
}
