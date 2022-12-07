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
			discord.EmbedError(s, i, discord.EmbedErrInvalidInteraction)
		} else {
			data := i.MessageComponentData()
			args := strings.Split(data.Values[0], "_")
			user := i.Member.User
			filler := event.FillerList[user.ID]
			room := event.RoomList[args[0]][args[1]]
			d, _ := strconv.Atoi(args[2])

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
				Data: &discordgo.InteractionResponseData{},
			})
			if err != nil {
				event.ErrExit(err)
			}

			// Check if deselectAll is selected.
			var deleteAll bool
			for _, v := range data.Values {
				if strings.HasSuffix(v, "deselect") {
					deleteAll = true
				}
			}

			if deleteAll {
				// If only default is selected, deschedules all hours for the day.
				for j := d * 24; j < (d+1)*24 && j < room.EventLength; j++ {
					if shiftIdx, has := event.HasShift(room, user.ID, j); has {
						room.Schedule[j][shiftIdx] = room.Schedule[j][len(room.Schedule[j])-1]
						room.Schedule[j] = room.Schedule[j][:len(room.Schedule[j])-1]
					}
				}
			} else {
				// Else schedules and deshedules user based on selection.
				shifts := []int{}
				for _, v := range data.Values {
					args := strings.Split(v, "_")
					hour := args[3]
					h, err := strconv.Atoi(hour)
					if err != nil {
						log.Fatal(err)
					}
					shifts = append(shifts, h)
				}

				for j := d * 24; j < (d+1)*24 && j < room.EventLength; j++ {
					if shiftIdx, has := event.HasShift(room, user.ID, j); has {
						if !inShift(shifts, j) {
							room.Schedule[j][shiftIdx] = room.Schedule[j][len(room.Schedule[j])-1]
							room.Schedule[j] = room.Schedule[j][:len(room.Schedule[j])-1]
						}
					} else {
						if inShift(shifts, j) {
							room.Schedule[j] = append(room.Schedule[j], filler.User.ID)
						}
					}
				}
			}

			embed := scheduleEmbeds(s, room, filler.Offset, d)
			compo := event.ScheduleDayComponent(room, room.EventLength/24+1, filler.Offset)

			_, err = s.FollowupMessageEdit(i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
				Embeds:     &embed,
				Components: &compo,
			})

			if err != nil {
				event.ErrExit(err)
			}

			// Log scheduling activities.
			log.Printf("%v has updated their day %v schedule for room '%v' in guild '%v'.", user.String(), d+1, room.Name, i.GuildID)

			// Backup room data.
			room.Backup()
		}
	}
}

func inShift(shifts []int, num int) bool {
	for _, hour := range shifts {
		if hour == num {
			return true
		}
	}
	return false
}
