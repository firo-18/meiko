package component

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/event"
)

func init() {
	List["hour-select"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Message.Interaction.User.String() != i.Member.User.String() {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Error",
							Description: "This interaction is intended for the original user only.",
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})

			if err != nil {
				log.Fatalln("interaction-respond:", err)
			}
		} else {
			data := i.MessageComponentData()

			args := strings.Split(data.Values[0], "_")
			key := args[0]
			room := event.RoomList[key]
			userID := i.Member.User.ID

			hourIdx, _ := strconv.Atoi(args[1])
			endIdx := hourIdx + 24
			if endIdx > room.EventLength {
				endIdx = room.EventLength
			}

			shifts := []int{}

			for _, v := range data.Values {
				if _, idx, ok := strings.Cut(v, "_"); ok {
					hour, err := strconv.Atoi(idx)
					if err != nil {
						log.Fatal(err)
					}
					shifts = append(shifts, hour)
				}
			}

			for j := hourIdx; j < endIdx; j++ {
				if shiftIdx, has := hasShift(userID, key, j); has {
					if !inShift(shifts, j) {
						room.Schedule[j][shiftIdx] = room.Schedule[j][len(room.Schedule[j])-1]
						room.Schedule[j] = room.Schedule[j][:len(room.Schedule[j])-1]

					}
				} else {
					if inShift(shifts, j) {
						fillerIdx, _ := findFiller(key, userID)
						room.Schedule[j] = append(room.Schedule[j], &room.Fillers[fillerIdx])
					}
				}
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					// Embeds:     scheduleEmbeds(s, key, day),
					Components: event.ScheduleDayComponent(room, room.EventLength/24+1),
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func findFiller(key, userID string) (int, bool) {
	room := event.RoomList[key]
	for i, v := range room.Fillers {
		if v.User.ID == userID {
			return i, true
		}
	}
	return 0, false
}

func inShift(shifts []int, num int) bool {
	for _, hour := range shifts {
		if hour == num {
			return true
		}
	}
	return false
}
