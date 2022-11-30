package event

// import (
// 	"fmt"
// 	"log"
// 	"regexp"
// 	"strings"
// 	"time"

// 	"github.com/bwmarrin/discordgo"
// 	"github.com/firo-18/meiko/client/discord"
// )

// func init() {
// 	List["shift"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		switch i.Type {
// 		case discordgo.InteractionApplicationCommand:
// 			data := i.ApplicationCommandData()

// 			key := data.Options[0].StringValue()

// 			if room, ok := RoomList[key]; !ok {
// 				discord.EmbedError(s, i, discord.EmbedErrorRoom404)
// 			} else {
// 				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 					Type: discordgo.InteractionResponseChannelMessageWithSource,
// 					Data: &discordgo.InteractionResponseData{
// 						Embeds: []*discordgo.MessageEmbed{
// 							{
// 								Title:     "Schedule for room - " + room.Name,
// 								Color:     discord.EmbedColor,
// 								Timestamp: discord.EmbedTimestamp,
// 								Footer:    discord.EmbedFooter(s),
// 								Fields:    shiftEmbedFields(s, i, key),
// 							},
// 						},
// 						Flags: discordgo.MessageFlagsEphemeral,
// 					},
// 				})
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 			}

// 		// Autocomplete
// 		case discordgo.InteractionApplicationCommandAutocomplete:
// 			data := i.ApplicationCommandData()
// 			choices := []*discordgo.ApplicationCommandOptionChoice{}
// 			choice := data.Options[0].StringValue()

// 			for _, v := range RoomList {
// 				_, ok := v.FindFiller(i.Member.User.ID)
// 				if v.Server == i.GuildID && ok {
// 					if ok, _ := regexp.MatchString("(?i)"+choice, v.Name); ok {
// 						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
// 							Name:  v.Name,
// 							Value: v.Key,
// 						})
// 					}
// 				}
// 			}

// 			// Max number of choice is 25.
// 			if len(choices) > 25 {
// 				choices = choices[:25]
// 			}

// 			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
// 				Data: &discordgo.InteractionResponseData{
// 					Choices: choices,
// 				},
// 			})
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	}
// }

// func shiftEmbedFields(s *discordgo.Session, i *discordgo.InteractionCreate, key string) []*discordgo.MessageEmbedField {
// 	fields := []*discordgo.MessageEmbedField{}

// 	room := RoomList[key]
// 	days := len(room.Schedule)/24 + 1

// 	for d := 0; d < days; d++ {
// 		values := []string{}
// 		for e := d * 24; e < (d+1)*24 && e < len(room.Schedule); e++ {
// 			if HasShift(i.Member.User.ID, key, e) {
// 				values = append(values, time.UnixMilli(room.Event.Start).Add(time.Hour*time.Duration(e)).Local().Format(OptionDateFormat))
// 			}
// 		}
// 		fields = append(fields, &discordgo.MessageEmbedField{
// 			Name:  fmt.Sprint("Day ", d+1),
// 			Value: discord.StyleFieldValues(strings.Join(values, ", ")),
// 		})
// 	}

// 	return fields
// }
