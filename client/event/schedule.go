package event

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

func init() {
	List["schedule"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			key := data.Options[0].StringValue()
			user := i.Member.User

			args := strings.Split(key, "_")

			// Check if room exist and return error if not.
			if room, ok := RoomList[args[0]][args[1]]; !ok {
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
									Fields:    discord.RoomInfoFields(room),
								},
							},
							Components: discord.DayScheduleComponents(room, filler, data.Name),
							Flags:      discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						LogError(err, data.Name)
					}

					time.Sleep(time.Minute * 5)

					_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Components: &[]discordgo.MessageComponent{},
					})

					if err != nil {
						LogError(err, data.Name)
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

	for _, room := range RoomList[i.GuildID] {
		if ok, _ := regexp.MatchString("(?i)"+choice, room.Name); ok {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  room.Name,
				Value: room.Key,
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
		LogError(err, data.Name)
	}
}

func HasShift(fillers []string, userID string) (int, bool) {
	for i, filler := range fillers {
		if filler == userID {
			return i, true
		}
	}
	return 0, false
}
