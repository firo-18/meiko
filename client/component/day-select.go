package component

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/client/event"
	"github.com/firo-18/meiko/schema"
)

func init() {
	List["day-select"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Message.Interaction.User.String() != i.Member.User.String() {
			discord.EmbedError(s, i, discord.EmbedErrInvalidInteraction)
		} else {
			data := i.MessageComponentData()
			args := strings.Split(data.Values[0], "_")
			key := args[0]
			day := args[1]
			d, err := strconv.Atoi(day)
			if err != nil {
				log.Fatal(err)
			}

			filler := event.FillerList[i.Member.User.ID]
			room := event.RoomList[key]

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     scheduleEmbeds(s, room, d, filler.Offset),
					Components: scheduleComponent(filler, room, d),
				},
			})
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func scheduleEmbeds(s *discordgo.Session, room *schema.Room, d, offset int) []*discordgo.MessageEmbed {
	embeds := []*discordgo.MessageEmbed{
		{
			Title:     fmt.Sprint("[Day ", d+1, "] Room - ", room.Name),
			Color:     discord.EmbedColor,
			Timestamp: discord.EmbedTimestamp,
			Footer:    discord.EmbedFooter(s),
			Fields:    scheduleEmbedFields(room, d, offset),
		},
	}

	return embeds
}

func scheduleEmbedFields(room *schema.Room, d, offset int) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}

	startTime := time.UnixMilli(room.Event.Start)

	for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
		eventTime := startTime.Add(time.Hour * time.Duration(h))
		if eventTime.Before(time.Now()) {
			continue
		}

		fillers := make([]string, len(room.Schedule[h]))
		for j, v := range room.Schedule[h] {
			fillers[j] = v.User.Username
		}
		value := strings.Join(fillers, ", ")
		if value == "" {
			value = "-"
		}
		timeOutput := eventTime.Add(time.Hour * time.Duration(offset*-1)).UTC().Format(discord.TimeOutputFormat)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprint("Hour ", h, " - ", timeOutput),
			Value:  discord.StyleFieldValues(value),
			Inline: true,
		})
	}

	return fields
}

func scheduleComponent(filler *schema.Filler, room *schema.Room, d int) []discordgo.MessageComponent {
	maxOption := 1

	startTime := time.UnixMilli(room.Event.Start)

	for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
		if startTime.Add(time.Hour * time.Duration(h)).After(time.Now()) {
			maxOption++
		}
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "hour-select",
					Placeholder: "Select your available hours.",
					MaxValues:   maxOption,
					Options:     scheduleComponentMenuOption(filler, room, d),
					Disabled:    maxOption == 1,
				},
			},
		},
	}

	return components
}

func scheduleComponentMenuOption(filler *schema.Filler, room *schema.Room, d int) []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{
		{
			Label:       "Deselect All",
			Value:       fmt.Sprint(room.Key, "_", d, "_", d*24, "_", "deselect"),
			Description: "Useful for when deselecting all hours .",
			Default:     false,
		},
	}

	startTime := time.UnixMilli(room.Event.Start)

	for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
		// Skip time options that are already passed.
		eventTime := startTime.Add(time.Hour * time.Duration(h))
		if eventTime.Before(time.Now()) {
			continue
		}

		timeOutput := eventTime.Add(time.Hour * time.Duration(filler.Offset*-1)).UTC().Format(discord.TimeOutputFormat)

		_, shift := hasShift(filler.User.ID, room.Key, h)

		options = append(options, discordgo.SelectMenuOption{
			Label:       timeOutput,
			Value:       fmt.Sprint(room.Key, "_", d, "_", h),
			Description: "Your offset time. If different from local, update through /link.",
			Default:     shift,
		})
	}

	return options
}

func hasShift(userID, key string, hour int) (int, bool) {
	for i, filler := range event.RoomList[key].Schedule[hour] {
		if filler.User.ID == userID {
			return i, true
		}
	}
	return 0, false
}
