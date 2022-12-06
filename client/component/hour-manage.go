package component

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/bwmarrin/discordgo"
// 	"github.com/firo-18/meiko/client/discord"
// 	"github.com/firo-18/meiko/client/event"
// 	"github.com/firo-18/meiko/schema"
// )

// func init() {
// 	List["hour-manage"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		if i.Message.Interaction.User.String() != i.Member.User.String() {
// 			discord.EmbedError(s, i, discord.EmbedErrInvalidInteraction)
// 		} else {
// 			data := i.MessageComponentData()
// 			args := strings.Split(data.Values[0], "_")
// 			key := args[0]
// 			room := event.RoomList[key]
// 			days := room.EventLength/24 + 1
// 			filler := event.FillerList[i.Member.User.ID]
// 			if len(args) > 3 {
// 				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 					Type: discordgo.InteractionResponseUpdateMessage,
// 					Data: &discordgo.InteractionResponseData{
// 						Components: event.ScheduleDayComponent(room, days, filler.Offset),
// 					},
// 				})
// 				if err != nil {
// 					event.ErrExit(err)
// 				}
// 				return
// 			}
// 			d, err := strconv.Atoi(args[1])
// 			if err != nil {
// 				event.ErrExit(err)
// 			}
// 			h, err := strconv.Atoi(args[2])
// 			if err != nil {
// 				event.ErrExit(err)
// 			}

// 			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 				Type: discordgo.InteractionResponseUpdateMessage,
// 				Data: &discordgo.InteractionResponseData{
// 					Embeds:     scheduleFillerEmbeds(s, room, d, h),
// 					Components: scheduleFillerComponents(room, h),
// 				},
// 			})
// 			if err != nil {
// 				event.ErrExit(err)
// 			}
// 		}
// 	}
// }

// func scheduleFillerEmbeds(s *discordgo.Session, room *schema.Room, d, h int) []*discordgo.MessageEmbed {
// 	fields := []*discordgo.MessageEmbedField{}

// 	for _, filler := range room.Schedule[h] {
// 		fields = append(fields, &discordgo.MessageEmbedField{
// 			Name:  filler.User.Username,
// 			Value: discord.FieldStyle(fmt.Sprintf("%v (%.2f)", filler.ISV, filler.SkillValue)),
// 		})
// 	}

// 	embed := &discordgo.MessageEmbed{
// 		Title:     fmt.Sprintf("Day %v - Hour %v", d+1, h),
// 		Color:     discord.EmbedColor,
// 		Timestamp: discord.EmbedTimestamp,
// 		Footer:    discord.EmbedFooter(s),
// 		Fields:    fields,
// 	}

// 	embeds := []*discordgo.MessageEmbed{embed}

// 	return embeds
// }

// func scheduleFillerComponents(room *schema.Room, h int) []discordgo.MessageComponent {
// 	options := []discordgo.SelectMenuOption{
// 		{
// 			Label:       "Default",
// 			Value:       fmt.Sprint(room.Key, "_", h, "_", "default"),
// 			Description: "Useful for when removing all fillers.",
// 			Default:     true,
// 		},
// 	}

// 	for _, filler := range room.Schedule[h] {
// 		options = append(options, discordgo.SelectMenuOption{
// 			Label:       filler.User.Username,
// 			Description: fmt.Sprintf("%v (%v)", filler.ISV, filler.SkillValue),
// 			Value:       fmt.Sprint(room.Key, "_", h, "_", filler.User.ID),
// 			Default:     true,
// 		})
// 	}

// 	menu := discordgo.SelectMenu{
// 		CustomID:    "filler-manage",
// 		Placeholder: "You can only remove filler(s). Irreversible.",
// 		MaxValues:   len(room.Schedule[h]) + 1,
// 		Options:     options,
// 	}

// 	components := []discordgo.MessageComponent{
// 		discordgo.ActionsRow{
// 			Components: []discordgo.MessageComponent{
// 				menu,
// 			},
// 		},
// 	}

// 	return components
// }
