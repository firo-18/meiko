package event

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/ghost"
)

func init() {
	List["join"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			user := i.Member.User
			var lead, sum int64
			var key string

			for i, v := range data.Options {
				switch v.Name {
				case "room":
					key = data.Options[i].StringValue()
				case "lead":
					lead = data.Options[i].IntValue()
				case "sum":
					sum = data.Options[i].IntValue()
				}
			}

			// Check if room exist and return error if not.
			if room, ok := RoomList[key]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrorRoom404)
			} else {
				// Calculate skill multiplier from ISV.
				skillValue := (float64(sum-lead) * 0.002) + float64(lead)/100 + 1

				// Check if user is already in the fillers list. If existed, update the info instead.
				if i, ok := FindFiller(RoomList[key].Fillers, *user); ok {
					room.Fillers[i] = ghost.New(*user, skillValue)
				} else {
					// Add user info to room if not found.
					room.Fillers = append(room.Fillers, ghost.New(*user, skillValue))
				}

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Success",
								Description: "You have successfully join to run/fill for the event. Run /scheduling to schedule your availability.",
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
										Name:  "ISV",
										Value: discord.StyleFieldValues(lead, "/", sum),
									},
									{
										Name:  "Skill Multiplier Value",
										Value: discord.StyleFieldValues(skillValue),
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
