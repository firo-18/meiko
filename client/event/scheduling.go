package event

import (
	"log"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

const (
	HourPerPage     = 24
	MaxScheduleHour = 72
)

func init() {
	List["scheduling"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			// user := i.Member.User
			key := data.Options[0].StringValue()

			// Check if room exist and return error if not.
			if room, ok := RoomList[key]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrorRoom404)
			} else {
				if _, ok := FindFiller(room.Fillers, *i.Member.User); !ok {
					discord.EmbedError(s, i, discord.EmbedErrorRoomNotJoined)
					return
				}

				if time.Now().UnixMilli() > room.Event.End {
					discord.EmbedError(s, i, discord.EmbedErrorRoomEnded)
					return
				}

				days := room.EventLength/24 + 1

				// Send select menu for scheduling if room exists.
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Components: ScheduleDayComponent(room, days),
						Flags:      discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					log.Fatal(err)
				}

				time.Sleep(time.Minute * 10)

				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Components: &[]discordgo.MessageComponent{},
				})

				if err != nil {
					log.Fatalln("interaction-respond:", err)
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
