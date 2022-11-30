package event

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

func init() {
	List["join"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			key := data.Options[0].StringValue()
			user := i.Member.User

			// Check if room exist and return error if not.
			if room, ok := RoomList[key]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrRoom404)
			} else {
				// Check to see if user is a filler.
				if filler, ok := FillerList[user.ID]; !ok {
					discord.EmbedError(s, i, discord.EmbedErrInvalidFiller)
				} else {
					// Add filler address to room's filler pool.
					room.FillerList[user.ID] = &filler

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:       "Success",
									Description: "You have successfully join the room. Run /schedule to select your availability.",
									Color:       discord.EmbedColor,
									Timestamp:   discord.EmbedTimestamp,
									Footer:      discord.EmbedFooter(s),
									Fields: []*discordgo.MessageEmbedField{
										{
											Name:  "Room",
											Value: discord.StyleFieldValues(room.Name),
										},
										{
											Name:  "Event",
											Value: discord.StyleFieldValues(RoomList[key].Event.Name),
										},
										{
											Name:  "Start Time",
											Value: discord.StyleFieldValues("<t:", room.Event.Start/1000, ">"),
										},
										{
											Name:  "Length",
											Value: discord.StyleFieldValues(room.EventLength, " hours"),
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
			}

		// Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			data := i.ApplicationCommandData()
			choices := []*discordgo.ApplicationCommandOptionChoice{}
			var choice string

			for i, v := range data.Options {
				if v.Name == "room" {
					choice = data.Options[i].StringValue()
				}
			}

			for _, v := range RoomList {
				if v.Server == i.GuildID {
					if ok, _ := regexp.MatchString("(?i)"+choice, v.Name); ok {
						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
							Name:  v.Name,
							Value: v.Key,
						})
					}
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
