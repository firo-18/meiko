package event

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/schema"
)

const (
	HourPerPage     = 24
	MaxScheduleHour = 72
)

func init() {
	List["schedule"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			key := data.Options[0].StringValue()
			user := i.Member.User

			// Check if room exist and return error if not.
			if room, ok := RoomList[key]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrRoom404)
			} else {
				// Check if user is a filler.
				if filler, ok := FillerList[user.ID]; !ok {
					discord.EmbedError(s, i, discord.EmbedErrInvalidFiller)
				} else {
					if time.Now().UnixMilli() > room.Event.End {
						discord.EmbedError(s, i, discord.EmbedErrRoomEnded)
						return
					}

					// Add filler to the room filler pool to keep track of who has offer support.
					room.FillerList[filler.User.ID] = &filler

					// Event length in days.
					daySum := room.EventLength/24 + 1

					// Send select menu for scheduling if room exists.
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:     fmt.Sprint("Room - ", room.Name),
									Color:     discord.EmbedColor,
									Timestamp: discord.EmbedTimestamp,
									Footer:    discord.EmbedFooter(s),
									Fields: []*discordgo.MessageEmbedField{
										{
											Name:  "Runner",
											Value: room.Name,
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
							Components: ScheduleDayComponent(room, daySum, filler.Offset),
							Flags:      discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						log.Fatal(err)
					}

					// Backup room data.
					room.Backup()

					time.Sleep(time.Minute * 5)

					_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Components: &[]discordgo.MessageComponent{},
					})

					if err != nil {
						log.Fatalln("interaction-respond:", err)
					}
				}
			}

		// Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			roomAutocomplete(s, i)
		}
	}
}

func roomAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	choice := data.Options[0].StringValue()

	for _, room := range RoomList {
		if i.GuildID == room.Server {
			if ok, _ := regexp.MatchString("(?i)"+choice, room.Name); ok {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  room.Name,
					Value: room.Key,
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

func ScheduleDayComponent(room *schema.Room, days, offset int) []discordgo.MessageComponent {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "day-select",
					Placeholder: "Select a day to schedule.",
					Options:     DaySelectMenu(room, days, offset),
				},
			},
		},
	}

	return components
}

func DaySelectMenu(room *schema.Room, days, offset int) []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{}

	for d := 0; d < days; d++ {

		// Find the last hour of the event day time.
		dayLastHour := (d+1)*24 - 1
		eventDayLastHour := time.UnixMilli(room.Event.Start).Add(time.Hour * time.Duration(dayLastHour))

		// Only add current day and beyond.
		if time.Now().Before(eventDayLastHour) {
			options = append(options, discordgo.SelectMenuOption{
				Label:       fmt.Sprint("Day ", d+1),
				Value:       fmt.Sprint(room.Key, "_", d),
				Description: "Start from " + time.UnixMilli(room.Event.Start).Add(time.Hour*24*time.Duration(d)).Add(time.Hour*time.Duration(offset)).UTC().Format(discord.TimeOutputFormat) + " offset time.",
			})
		}
	}
	return options
}
