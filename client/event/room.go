package event

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/room"
)

func init() {
	List["room"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			var name, event string

			switch data.Options[0].Name {
			case "name":
				name, event = data.Options[0].StringValue(), data.Options[1].StringValue()
			case "event":
				event, name = data.Options[0].StringValue(), data.Options[1].StringValue()
			}

			name = strings.ToUpper(name)

			if _, ok := RoomList[name]; ok {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Error",
								Description: discord.StyleFieldValues("Room name '", name, "' is already existed. Choose a different name."),
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
				idx, err := strconv.Atoi(event)
				if err != nil {
					log.Fatal(err)
				}

				idx-- // EventID starts at 1.
				if idx > 37 {
					idx-- // RMD shenanigan. Missing 1 event basically, so EventID jumps from 37 to 39.
				}

				RoomList[name] = room.New(EventList[idx], *i.Member.User)

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Success",
								Description: discord.StyleFieldValues("Room '", name, "' has been created."),
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
			}

		// Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			data := i.ApplicationCommandData()
			choices := []*discordgo.ApplicationCommandOptionChoice{}
			var choice string

			switch data.Options[0].Name {
			case "name":
				choice = data.Options[1].StringValue()
			case "event":
				choice = data.Options[0].StringValue()
			}

			for _, v := range EventList {
				if v.End > time.Now().UnixMilli() {
					if ok, _ := regexp.MatchString("(?i)"+choice, v.Name); ok {
						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
							Name:  v.Name,
							Value: strconv.Itoa(v.ID),
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
