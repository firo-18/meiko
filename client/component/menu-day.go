package component

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/client/event"
)

func init() {
	List["menu-day"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			filler := event.FillerList[i.Member.User.ID]
			day := args[2]
			d, err := strconv.Atoi(day)
			if err != nil {
				event.ErrExit(err)
			}
			components := []discordgo.MessageComponent{}

			switch i.Message.Interaction.Name {
			case "view":
				components = discord.DayScheduleComponents(room, filler, i.Message.Interaction.Name)
			case "schedule", "manage":
				components = discord.HourScheduleComponents(room, filler, d, i.Message.Interaction.Name)
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     discord.DayScheduleEmbeds(s, event.FillerList, room, filler, d, i.Message.Interaction.Name),
					Components: components,
				},
			})
			if err != nil {
				event.ErrExit(err)
			}
		}
	}
}
