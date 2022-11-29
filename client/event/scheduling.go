package event

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/room"
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

func ScheduleDayComponent(room *room.Room, days int) []discordgo.MessageComponent {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "day-select",
					Placeholder: "Select a day to schedule.",
					Options:     DaySelectMenu(room, days),
				},
			},
		},
	}

	return components
}

func DaySelectMenu(room *room.Room, days int) []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{}

	for i := 0; i < days; i++ {
		options = append(options, discordgo.SelectMenuOption{
			Label:       fmt.Sprint("Day ", i+1),
			Value:       fmt.Sprint(room.Key, "_", i+1),
			Description: time.UnixMilli(room.Event.Start).Add(time.Hour * 24 * time.Duration(i)).Local().Format(OptionDateFormat),
		})
	}
	return options
}
