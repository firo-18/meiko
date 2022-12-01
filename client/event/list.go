package event

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/schema"
)

func init() {
	List["list"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			var key string
			var deleteRoom bool

			for _, option := range data.Options {
				switch option.Name {
				case "room":
					key = option.StringValue()
				case "delete-room":
					deleteRoom = true
				}
			}
			if room, ok := RoomList[key]; !ok {
				discord.EmbedError(s, i, discord.EmbedErrRoom404)
			} else {
				if deleteRoom {
					user := i.Member.User
					if user.ID != room.Owner.ID {
						discord.EmbedError(s, i, discord.EmbedErrInvalidOwner)
						return
					}

					// Archives as json before deleting room.
					err := room.Archive()
					if err != nil {
						log.Fatal(err)
					}
					delete(RoomList, key)
					err = os.Remove(schema.PathRoomDB + room.Key + ".gob")
					if err != nil {
						log.Fatal(err)
					}

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
						log.Fatal(err)
					}

					log.Printf("%v has deleted the room '%v' from guild '%v'.", user.String(), room.Name, room.Server)
				} else {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:     "Viewing Room Info - " + room.Name,
									Color:     discord.EmbedColor,
									Timestamp: discord.EmbedTimestamp,
									Footer:    discord.EmbedFooter(s),
									Fields:    discord.RoomInfoFields(room),
								},
							},
							Flags: discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						log.Fatal(err)
					}
				}
			}

		// Room Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			roomAutocomplete(s, i)
		}
	}
}
