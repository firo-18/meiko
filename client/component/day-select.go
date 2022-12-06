package component

import (
	"fmt"
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
				event.ErrExit(err)
			}

			filler := event.FillerList[i.Member.User.ID]
			room := event.RoomList[key]

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

func scheduleEmbeds(s *discordgo.Session, room *schema.Room, d, offset int) []*discordgo.MessageEmbed {
	local := time.UnixMilli(room.Event.Start).Add(time.Hour * time.Duration(d))
	off := local.Add(time.Hour * time.Duration(offset)).UTC().Format(discord.TimeOutputFormat)

	embeds := []*discordgo.MessageEmbed{
		{
			Title:       fmt.Sprint("[Day ", d+1, "] Room - ", room.Name),
			Description: fmt.Sprintf("Local: <t:%v:f> | Offset: %v. If the two times are different, update your offset with /link.", local.Unix(), off),
			Color:       discord.EmbedColor,
			Timestamp:   discord.EmbedTimestamp,
			Footer:      discord.EmbedFooter(s),
			Fields:      scheduleEmbedFields(room, d, offset),
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
		timeOutput := eventTime.Add(time.Hour * time.Duration(offset)).UTC().Format(discord.TimeOutputFormat)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprint("Hour ", h, " - ", timeOutput),
			Value:  discord.FieldStyle(value),
			Inline: true,
		})
	}

	return fields
}

func scheduleComponent(i *discordgo.InteractionCreate, filler *schema.Filler, room *schema.Room, d int) []discordgo.MessageComponent {

	menu := discordgo.SelectMenu{
		CustomID:    "hour-select",
		Placeholder: "Select your available hours.",
		Options:     scheduleComponentMenuOption(i, filler, room, d),
	}
	switch i.Message.Interaction.Name {
	case "schedule":
		maxOption := 1
		startTime := time.UnixMilli(room.Event.Start)

		for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
			if startTime.Add(time.Hour * time.Duration(h)).After(time.Now()) {
				maxOption++
			}
		}
		menu.MaxValues = maxOption
		menu.Disabled = maxOption == 1
	case "manage":
		menu.CustomID = "hour-manage"
		menu.MaxValues = 1
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				menu,
			},
		},
	}

	return components
}

func scheduleComponentMenuOption(i *discordgo.InteractionCreate, filler *schema.Filler, room *schema.Room, d int) []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{}
	switch i.Message.Interaction.Name {
	case "schedule":
		options = append(options, discordgo.SelectMenuOption{
			Label:       "Deselect All",
			Value:       fmt.Sprint(room.Key, "_", d, "_", d*24, "_", "deselect"),
			Description: "Useful for when deselecting all hours.",
			Default:     false,
		})
	case "manage":
		options = append(options, discordgo.SelectMenuOption{
			Label:       "Back",
			Value:       fmt.Sprint(room.Key, "_", d, "_", d*24, "_", "back"),
			Description: "Go back to the previous options.",
			Default:     false,
		})
	}

	startTime := time.UnixMilli(room.Event.Start)

	for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
		// Skip time options that are already passed.
		eventTime := startTime.Add(time.Hour * time.Duration(h))
		if eventTime.Before(time.Now()) {
			continue
		}

		timeOutput := eventTime.Add(time.Hour * time.Duration(filler.Offset)).UTC().Format(discord.TimeOutputFormat)

		switch i.Message.Interaction.Name {
		case "schedule":
			_, shift := event.HasShift(filler.User.ID, room.Key, h)

			options = append(options, discordgo.SelectMenuOption{
				Label:       timeOutput,
				Value:       fmt.Sprint(room.Key, "_", d, "_", h),
				Description: fmt.Sprint("Event Hour: ", h),
				Default:     shift,
			})
		case "manage":
			if len(room.Schedule[h]) > 0 {
				options = append(options, discordgo.SelectMenuOption{
					Label:       timeOutput,
					Value:       fmt.Sprint(room.Key, "_", d, "_", h),
					Description: fmt.Sprint("Event Hour: ", h),
				})
			}
		}

	}

	return options
}
