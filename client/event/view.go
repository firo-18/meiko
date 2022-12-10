package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

func init() {
	List["view"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
