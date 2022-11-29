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
			key := i.GuildID + " - " + name

			if _, ok := RoomList[key]; ok {
				discord.EmbedError(s, i, discord.EmbedErrorRoomNameDuplicated)
			} else {
				idx, err := strconv.Atoi(event)
				if err != nil {
					log.Fatal(err)
				}

				idx-- // EventID starts at 1.
				if idx > 37 {
					idx-- // RMD shenanigan. Missing 1 event basically, so EventID jumps from 37 to 39.
				}

				RoomList[key] = room.New(EventList[idx], *i.Member.User, name, i.GuildID)

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:     "Room Created",
								Color:     discord.EmbedColor,
								Timestamp: discord.EmbedTimestamp,
								Footer:    discord.EmbedFooter(s),
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:  "Name",
										Value: name,
									},
									{
										Name:  "Event",
										Value: EventList[idx].Name,
									},
									{
										Name:  "Server",
										Value: i.GuildID,
									},
									{
										Name:  "Created By",
										Value: i.Member.User.Username,
									},
								},
							},
						},
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
