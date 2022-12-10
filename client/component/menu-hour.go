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
	List["menu-hour"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			user := i.Member.User
			filler := event.FillerList[user.ID]
			d, _ := strconv.Atoi(args[2])

			switch i.Message.Interaction.Name {
			case "schedule":
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
					Data: &discordgo.InteractionResponseData{},
				})
				if err != nil {
					event.LogError(err, i.Message.Interaction.Name)
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
					for h := d * 24; h < (d+1)*24 && h < room.EventLength; h++ {
						if shiftIdx, has := event.HasShift(room.Schedule[h], user.ID); has {
							room.Schedule[h][shiftIdx] = room.Schedule[h][len(room.Schedule[h])-1]
							room.Schedule[h] = room.Schedule[h][:len(room.Schedule[h])-1]
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
						if shiftIdx, has := event.HasShift(room.Schedule[j], user.ID); has {
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

				embed := discord.DayScheduleEmbeds(s, event.FillerList, room, filler, d, i.Message.Interaction.Name)
				compo := discord.DayScheduleComponents(room, filler, i.Message.Interaction.Name)

				_, err = s.FollowupMessageEdit(i.Interaction, i.Message.ID, &discordgo.WebhookEdit{
					Embeds:     &embed,
					Components: &compo,
				})

				if err != nil {
					event.LogError(err, i.Message.Interaction.Name)
				}

				// Log scheduling activities.
				log.Printf("%v has updated their day %v schedule for room '%v' in guild '%v'.", user.String(), d+1, room.Name, i.GuildID)

				// Backup room data.
				room.Backup()

			case "manage":
				if len(args) > 4 {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseUpdateMessage,
						Data: &discordgo.InteractionResponseData{
							Components: discord.DayScheduleComponents(room, filler, i.Message.Interaction.Name),
						},
					})
					if err != nil {
						event.LogError(err, i.Message.Interaction.Name)
					}
					return
				}
				d, err := strconv.Atoi(args[2])
				if err != nil {
					event.LogError(err, i.Message.Interaction.Name)
				}
				h, err := strconv.Atoi(args[3])
				if err != nil {
					event.LogError(err, i.Message.Interaction.Name)
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds:     discord.FillerScheduleEmbeds(s, event.FillerList, room, d, h),
						Components: discord.FillerScheduleComponents(event.FillerList, room, h),
					},
				})
				if err != nil {
					event.LogError(err, i.Message.Interaction.Name)
				}
			}

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
