package component

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/client/event"
)

func init() {
	List["menu-fillers"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Message.Interaction.User.String() != i.Member.User.String() {
			discord.EmbedError(s, i, discord.EmbedErrInvalidInteraction)
		} else {
			data := i.MessageComponentData()
			args := strings.Split(data.Values[0], "_")
			room := event.RoomList[args[0]][args[1]]
			if room == nil {
				discord.EmbedError(s, i, discord.EmbedErrRoomDeleted)
				return
			}

			h, err := strconv.Atoi(args[2])
			if err != nil {
				event.ErrExit(err)
			}

			room.Schedule[h] = []string{}

			for _, v := range data.Values {
				args := strings.Split(v, "_")
				if args[3] != "default" {
					room.Schedule[h] = append(room.Schedule[h], args[3])
				}
			}

			d := h / 24
			filler := event.FillerList[i.Member.User.ID]

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     discord.DayScheduleEmbeds(s, event.FillerList, room, filler, d, i.Message.Interaction.Name),
					Components: discord.HourScheduleComponents(room, filler, d, i.Message.Interaction.Name),
				},
			})
			if err != nil {
				event.ErrExit(err)
			}

			// Log scheduling activities.
			log.Printf("%v has updated fillers for shift %v for room '%v' in guild '%v'.", i.Member.User.String(), h, room.Name, i.GuildID)

			// Backup room data.
			room.Backup()
		}
	}
}
