package event

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/schema"
)

type ChannelMutex struct {
	mu    sync.Mutex
	quits map[string]chan bool
}

var (
	chMu = ChannelMutex{
		quits: make(map[string]chan bool),
	}
	orderCC   = []int{3, 5, 4, 1, 2}
	orderEnvy = []int{5, 4, 1, 2, 3}
)

func init() {
	List["session"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			var key, cmd string
			var role *discordgo.Role

			for j, v := range data.Options {
				switch v.Name {
				case "room":
					key = data.Options[j].StringValue()
				case "state":
					cmd = data.Options[j].StringValue()
				case "role":
					role = data.Options[j].RoleValue(s, i.GuildID)
				}
			}

			room, ok := RoomList[key]
			if !ok {
				discord.EmbedError(s, i, discord.EmbedErrRoom404)
				return
			}

			user := i.Member.User

			if room.Owner.ID != user.ID {
				discord.EmbedError(s, i, discord.EmbedErrInvalidOwner)
				return
			}

			chMu.mu.Lock()
			chMu.quits[key] = make(chan bool, 2)
			chMu.mu.Unlock()

			if cmd == "begin" {
				go func() {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title:       fmt.Sprintf("[Room: %v] Tiering Session Has Begun", room.Name),
									Description: "Begin sending out check in message roughly 15 minutes before shift starts.",
									Color:       discord.EmbedColor,
									Timestamp:   discord.EmbedTimestamp,
									Footer:      discord.EmbedFooter(s),
								},
							},
						},
					})
					if err != nil {
						ErrExit(err)
					}

					// Log session activities
					log.Printf("%v started a session for room '%v' in guild %v.", user.String(), room.Name, i.GuildID)

					first := true
					h := 0
					for h < len(room.Schedule) {
						chMu.mu.Lock()
						select {
						case <-chMu.quits[key]:
							chMu.mu.Unlock()
							return
						default:
							chMu.mu.Unlock()
							currTime := time.Now()
							eventStartTime := time.UnixMilli(room.Event.Start)

							nextHourTime := eventStartTime.Add(time.Hour * time.Duration(h))
							if currTime.After(nextHourTime) {
								h++
								continue
							}
							if nextHourTime.Sub(currTime) <= time.Duration(time.Minute*15) {
								fillers := room.Schedule[h]
								sort.SliceStable(fillers, func(i, j int) bool {
									return fillers[i].SkillValue > fillers[j].SkillValue
								})

								roomOrder := make([]*schema.Filler, len(fillers))
								order := []int{}
								switch room.Event.Type {
								case "marathon":
									order = orderEnvy
								case "cheerful_carnival":
									order = orderCC
								}
								idx := 0
								for _, v := range order {
									if v > len(fillers) {
										continue
									}
									roomOrder[v-1] = fillers[idx]
									idx++
								}

								if err != nil {
									ErrExit(err)
								}

								roomOrderMention := []string{}
								for _, f := range roomOrder {
									if _, ok := HasShift(f.User.ID, key, h-1); !ok || first {
										roomOrderMention = append(roomOrderMention, f.User.Mention())
									}
								}

								if role != nil {
									roomOrderMention = append(roomOrderMention, role.Mention())
								}

								embed := discord.NewEmbed(s)
								embed.Title = fmt.Sprintf("[Room: %v] Shift Check In: <t:%v:R>", room.Name, nextHourTime.Unix())
								embed.Description = fmt.Sprintf("Event: %v - %v", room.Event.Name, room.Event.Type)
								embed.Fields = fillerOrderFields(roomOrder, role)

								embeds := []*discordgo.MessageEmbed{embed}

								_, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
									Content: strings.Join(roomOrderMention, " | "),
									Embeds:  embeds,
								})
								if err != nil {
									s.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
										Title:     fmt.Sprintf("[Room: %v] Tiering Session Has Ended Dues To Error", room.Name),
										Color:     discord.EmbedColor,
										Timestamp: discord.EmbedTimestamp,
										Footer:    discord.EmbedFooter(s),
									})
									ErrExit(err)
								}
								h++
								first = false
							}
							time.Sleep(time.Minute)
						}
					}
				}()
			} else {
				// Send defer respond to interaction
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{},
				})
				if err != nil {
					ErrExit(err)
				}

				// End goroutine if cmd is 'stop'
				chMu.mu.Lock()
				defer chMu.mu.Unlock()
				chMu.quits[key] <- true

				// Send follow up session ending message.
				embeds := []*discordgo.MessageEmbed{
					{
						Title:       fmt.Sprintf("[Room: %v] Tiering Session Has Ended", room.Name),
						Description: "Stop sending out check in messages.",
						Color:       discord.EmbedColor,
						Timestamp:   discord.EmbedTimestamp,
						Footer:      discord.EmbedFooter(s),
					},
				}
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Embeds: &embeds,
				})
				if err != nil {
					ErrExit(err)
				}

				// Log session activities
				log.Printf("%v ended the session for room '%v' in guild %v.", user.String(), room.Name, i.GuildID)
			}

		// Room Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			roomAutocomplete(s, i)
		}
	}
}

func fillerOrderFields(fillers []*schema.Filler, role *discordgo.Role) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}

	for i, v := range fillers {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprint("Player ", i+1),
			Value: discord.FieldStyle(v.User.Mention(), " - ", v.SkillValue),
		})
	}

	if len(fields) < 5 && role != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Asking For Additional Filler(s)",
			Value: discord.FieldStyle(fmt.Sprintf("%v: (+%v)", role.Mention(), 5-len(fields))),
		})
	}

	return fields
}
