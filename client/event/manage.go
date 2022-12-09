package event

import (
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

func init() {
	List["manage"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			var key string
			var manager discordgo.User
			var deleteRoom bool
			user := i.Member.User

			for _, option := range data.Options {
				switch option.Name {
				case "room":
					key = option.StringValue()
				case "manager":
					manager = *option.UserValue(s)
				case "delete-room":
					deleteRoom = true
				}
			}

			args := strings.Split(key, "_")
			if room, ok := RoomList[args[0]][args[1]]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrRoom404)
			} else {
				if user.ID != room.Owner.ID && user.ID != room.Manager.ID {
					discord.EmbedError(s, i, discord.EmbedErrInvalidOwner)
					return
				}
				if deleteRoom {
					// Backup as json before deleting room.
					err := room.Archive()
					if err != nil {
						ErrExit(err)
					}

					delete(RoomList[args[0]], args[1])

					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:     "Room Deleted - " + room.Name,
									Color:     discord.EmbedColor,
									Timestamp: discord.EmbedTimestamp,
									Footer:    discord.EmbedFooter(s),
									Fields:    discord.RoomInfoFields(room),
								},
							},
						},
					})
					if err != nil {
						ErrExit(err)
					}

					log.Printf("%v has deleted the room '%v' from guild '%v'.", user.String(), room.Name, room.Guild)
				} else {
					if time.Now().UnixMilli() > room.Event.End {
						discord.EmbedError(s, i, discord.EmbedErrRoomEnded)
						return
					}
					// Add/update manager if a manager is selected.
					if manager.ID != "" {
						room.Manager = manager
						// Log manager's changes.
						log.Printf("%v changed the manager of room '%v' in guild %v to %v.", user.String(), room.Name, i.GuildID, manager.String())
					}

					filler, ok := FillerList[user.ID]
					if !ok {
						discord.EmbedError(s, i, discord.EmbedErrInvalidFiller)
						return
					}

					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:     "Managing Room - " + room.Name,
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
						ErrExit(err)
					}

					// Backup room data.
					room.Backup()

					time.Sleep(time.Minute * 5)

					_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Components: &[]discordgo.MessageComponent{},
					})

					if err != nil {
						ErrExit(err)
					}

				}
			}

		// Room Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			roomAutocomplete(s, i)
		}
	}
}
