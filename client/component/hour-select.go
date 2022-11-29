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
	List["hour-select"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Message.Interaction.User.String() != i.Member.User.String() {
			discord.EmbedError(s, i, discord.EmbedErrorInvalidInteraction)
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

			// Check if default is selected.
			var isDefault bool
			if len(args) > 2 {
				isDefault = true
			}
			if isDefault && len(data.Values) == 1 {
				// If only default is selected, deschedules all hours for the day.
				for j := hourIdx; j < endIdx; j++ {
					if shiftIdx, has := hasShift(userID, key, j); has {
						room.Schedule[j][shiftIdx] = room.Schedule[j][len(room.Schedule[j])-1]
						room.Schedule[j] = room.Schedule[j][:len(room.Schedule[j])-1]
					}
				}
			} else {
				// Else schedules and deshedules user based on selection.

				start := 0
				if isDefault {
					start = 1
				}

				shifts := []int{}
				for _, v := range data.Values[start:] {
					args := strings.Split(v, "_")
					idx := args[1]
					hour, err := strconv.Atoi(idx)
					if err != nil {
						log.Fatal(err)
					}
					shifts = append(shifts, hour)
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

			}

			dayIdx := hourIdx / 24

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     scheduleEmbeds(s, key, dayIdx),
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
