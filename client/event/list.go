package event

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

func init() {
	List["list"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			if room, ok := RoomList[data.Options[0].StringValue()]; !ok {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Error",
								Description: "Room not exist. Select a room from the list.",
								Color:       discord.EmbedColor,
								Timestamp:   discord.EmbedTimestamp,
								Footer:      discord.EmbedFooter(s),
							},
						},
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					log.Fatal(err)
				}
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:     "Room Info - " + data.Options[0].StringValue(),
								Color:     discord.EmbedColor,
								Timestamp: discord.EmbedTimestamp,
								Footer:    discord.EmbedFooter(s),
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:  "Event",
										Value: discord.StyleFieldValues(room.Event.Name),
									},
									{
										Name:  "Created By",
										Value: discord.StyleFieldValues(room.Creator.Username),
									},
									{
										Name:  "Runner",
										Value: discord.StyleFieldValues(room.Runner),
									},
									{
										Name:  "Fillers",
										Value: discord.StyleFieldValues(len(room.Fillers), " fillers signed up."),
									},
								},
							},
						},
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					log.Fatal(err)
				}
			}

		// Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			data := i.ApplicationCommandData()
			choices := []*discordgo.ApplicationCommandOptionChoice{}
			choice := data.Options[0].StringValue()

			for k := range RoomList {
				if ok, _ := regexp.MatchString("(?i)"+choice, k); ok {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  k,
						Value: k,
					})
				}
			}

			// Max number of choice is 25.
			if len(choices) > 25 {
				choices = choices[:25]
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
