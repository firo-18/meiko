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
)

func init() {
	List["day-select"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Message.Interaction.User.String() != i.Member.User.String() {
			discord.EmbedError(s, i, discord.EmbedErrorInvalidInteraction)
		} else {
			data := i.MessageComponentData()

			args := strings.Split(data.Values[0], "_")
			key := args[0]
			day := args[1]
			dayIdx, err := strconv.Atoi(day)
			if err != nil {
				log.Fatal(err)
			}
			dayIdx-- // Day 1 is 0 index

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     scheduleEmbeds(s, key, dayIdx),
					Components: scheduleComponent(key, i.Member.User.ID, dayIdx),
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func scheduleEmbeds(s *discordgo.Session, key string, dayIdx int) []*discordgo.MessageEmbed {
	embeds := []*discordgo.MessageEmbed{
		{
			Title:     fmt.Sprint("Schedule - Day ", dayIdx+1),
			Color:     discord.EmbedColor,
			Timestamp: discord.EmbedTimestamp,
			Footer:    discord.EmbedFooter(s),
			Fields:    scheduleEmbedFields(key, dayIdx),
		},
	}

	return embeds
}

func scheduleEmbedFields(key string, dayIdx int) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}

	schedule := event.RoomList[key].Schedule
	startTime := time.UnixMilli(event.RoomList[key].Event.Start)

	for i := dayIdx * 24; i < (dayIdx+1)*24 && i < len(schedule); i++ {
		fillers := make([]string, len(schedule[i]))
		for j, v := range schedule[i] {
			fillers[j] = v.User.Username
		}
		value := strings.Join(fillers, ", ")
		if value == "" {
			value = "No filler signed up."
		}
		timeOutput := startTime.Add(time.Hour * time.Duration(i)).Format(discord.TimeOutputFormat)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprint("Hour ", i, " - ", timeOutput),
			Value:  discord.StyleFieldValues(value),
			Inline: true,
		})
	}

	return fields
}

func scheduleComponent(key, userID string, dayIdx int) []discordgo.MessageComponent {
	maxOption := 0

	if dayIdx < len(event.RoomList[key].Schedule)/24 {
		maxOption = 25
	} else {
		maxOption = 7
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "hour-select",
					Placeholder: "Select your available hours.",
					MaxValues:   maxOption,
					Options:     scheduleComponentMenuOption(key, userID, dayIdx),
				},
			},
		},
	}

	return components
}

func scheduleComponentMenuOption(key, userID string, dayIdx int) []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{
		{
			Label:       "Default",
			Value:       fmt.Sprint(key, "_", dayIdx*24, "_", "default"),
			Description: "Keep this selected. Useful for when deselecting all hours.",
			Default:     true,
		},
	}

	schedule := event.RoomList[key].Schedule
	startTime := time.UnixMilli(event.RoomList[key].Event.Start)

	for i := dayIdx * 24; i < (dayIdx+1)*24 && i < len(schedule); i++ {
		timeOutput := startTime.Add(time.Hour * time.Duration(i)).Local().Format(discord.TimeOutputFormat)

		_, shift := hasShift(userID, key, i)

		options = append(options, discordgo.SelectMenuOption{
			Label:       timeOutput,
			Value:       fmt.Sprint(key, "_", i),
			Description: "Local time.",
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
