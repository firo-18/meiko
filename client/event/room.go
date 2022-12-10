package event

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/schema"
)

func init() {
	List["room"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()

			// Check whether user has linked filler account.
			if _, ok := FillerList[i.Member.User.ID]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrInvalidFiller)
				return
			}

			var roomName, eventName string
			var fill bool
			for i, v := range data.Options {
				switch v.Name {
				case "name":
					roomName = data.Options[i].StringValue()
				case "event":
					eventName = data.Options[i].StringValue()
				case "fill-all":
					fill = data.Options[i].BoolValue()
				}
			}

			// Check if event name exists, and return err if not.
			if event, ok := EventList[eventName]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrInvalidEvent)
			} else {
				roomName = strings.ToUpper(roomName)

				if _, ok := RoomList[i.GuildID][roomName]; ok {
					discord.EmbedError(s, i, discord.EmbedErrRoomNameDuplicated)
				} else {
					user := i.Member.User
					filler := FillerList[user.ID]

					if list, ok := RoomList[i.GuildID]; ok {
						list[roomName] = schema.NewRoom(i.GuildID, roomName, event, user)
					} else {
						rooms := make(map[string]*schema.Room)
						rooms[roomName] = schema.NewRoom(i.GuildID, roomName, event, user)
						RoomList[i.GuildID] = rooms
					}

					room := RoomList[i.GuildID][roomName]

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:       "Room Created - " + room.Name,
									Description: "You can run /schedule to alter your availability.",
									Color:       discord.EmbedColor,
									Timestamp:   discord.EmbedTimestamp,
									Footer:      discord.EmbedFooter(s),
									Fields:      discord.RoomInfoFields(room),
								},
							},
						},
					})
					if err != nil {
						LogError(err, data.Name)
					}

					// If fill-all is selected, add runner to all hour slots.
					if fill {
						for j := 0; j < len(room.Schedule); j++ {
							room.Schedule[j] = append(room.Schedule[j], filler.User.ID)
						}
					}

					// Log room creation activities.
					log.Printf("%v created a room named '%v' in guild %v.", user.String(), roomName, i.GuildID)

					// Backup room data.
					room.Backup()
				}
			}

		// Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			eventAutocomplete(s, i)
		}
	}
}

func eventAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	var choice string

	for j, v := range data.Options {
		if v.Name == "event" {
			choice = data.Options[j].StringValue()
		}
	}

	for _, v := range EventList {
		if v.End > time.Now().UnixMilli() {
			if ok, _ := regexp.MatchString("(?i)"+choice, v.Name); ok {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  v.Name,
					Value: v.Name,
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
		LogError(err, data.Name)
	}
}
