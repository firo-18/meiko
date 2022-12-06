package component

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/client/event"
	"github.com/firo-18/meiko/schema"
)

func init() {
	List["filler-manage"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Message.Interaction.User.String() != i.Member.User.String() {
			discord.EmbedError(s, i, discord.EmbedErrInvalidInteraction)
		} else {
			data := i.MessageComponentData()
			args := strings.Split(data.Values[0], "_")
			key := args[0]
			room := event.RoomList[key]

			h, err := strconv.Atoi(args[1])
			if err != nil {
				event.ErrExit(err)
			}

			room.Schedule[h] = []*schema.Filler{}

			if len(data.Values) > 1 || args[2] != "default" {
				for _, v := range data.Values {
					args := strings.Split(v, "_")
					filler := event.FillerList[args[2]]
					room.Schedule[h] = append(room.Schedule[h], filler)
				}
			}

			d := h / 24
			filler := event.FillerList[i.Member.User.ID]

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     scheduleEmbeds(s, room, d, filler.Offset),
					Components: scheduleComponent(i, filler, room, d),
				},
			})
			if err != nil {
				event.ErrExit(err)
			}
		}
	}
}
