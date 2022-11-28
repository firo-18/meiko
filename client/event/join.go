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
			var room string

			for i, v := range data.Options {
				switch v.Name {
				case "room":
					room = data.Options[i].StringValue()
				case "lead":
					lead = data.Options[i].IntValue()
				case "sum":
					sum = data.Options[i].IntValue()
				}
			}

			// Check if room exist and return error if not.
			if _, ok := RoomList[room]; !ok {
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
				// Calculate skill multiplier from ISV.
				skillValue := (float64(sum-lead) * 0.002) + float64(lead)/100 + 1

				// Check if user is already in the fillers list. If existed, update the info instead.
				if i, ok := findFiller(RoomList[room].Fillers, *user); ok {
					RoomList[room].Fillers[i] = ghost.New(*user, skillValue)
				} else {
					// Add user info to room if not found.
					RoomList[room].Fillers = append(RoomList[room].Fillers, ghost.New(*user, skillValue))
				}

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Success",
								Description: "You have successfully sign up to run/fill for the event.",
								Color:       discord.EmbedColor,
								Timestamp:   discord.EmbedTimestamp,
								Footer:      discord.EmbedFooter(s),
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:  "Room",
										Value: discord.StyleFieldValues(room),
									},
									{
										Name:  "Event",
										Value: discord.StyleFieldValues(RoomList[room].Event.Name),
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

// findFiller loops through the fillers slice to find a filler. If exists, it returns the index and true, otherwise, it returns 0 and false.
func findFiller(fillers []*ghost.Ghost, filler discordgo.User) (int, bool) {
	for i, v := range fillers {
		if v.User == filler {
			return i, true
		}
	}
	return 0, false
}
