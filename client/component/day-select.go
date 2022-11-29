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

			args := strings.Split(data.Values[0], " - ")
			key := args[0] + " - " + args[1]
			day := args[2]

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     scheduleEmbeds(s, key, day),
					Components: scheduleComponent(key, day, i.Member.User.ID),
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func scheduleEmbeds(s *discordgo.Session, key, day string) []*discordgo.MessageEmbed {
	embeds := []*discordgo.MessageEmbed{
		{
			Title:     "Schedule - Day " + day,
			Color:     discord.EmbedColor,
			Timestamp: discord.EmbedTimestamp,
			Footer:    discord.EmbedFooter(s),
			Fields:    scheduleEmbedFields(key, day),
		},
	}

	return embeds
}

func scheduleEmbedFields(key, day string) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}

	dayIdx, err := strconv.Atoi(day)
	if err != nil {
		log.Fatal(err)
	}
	dayIdx-- // Day 1 is 0 index

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

func scheduleComponent(key, day, userID string) []discordgo.MessageComponent {
	dayNum, err := strconv.Atoi(day)
	if err != nil {
		log.Fatal(err)
	}

	maxOption := 0

	if dayNum < len(event.RoomList[key].Schedule)/24+1 {
		maxOption = 24
	} else {
		maxOption = 6
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "hour-select",
					Placeholder: "Select your available hours.",
					MaxValues:   maxOption,
					Options:     scheduleComponentMenuOption(key, day, userID),
				},
			},
		},
	}

	return components
}

func scheduleComponentMenuOption(key, day, userID string) []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{}

	dayIdx, err := strconv.Atoi(day)
	if err != nil {
		log.Fatal(err)
	}
	dayIdx-- // Day 1 is 0 index

	schedule := event.RoomList[key].Schedule
	startTime := time.UnixMilli(event.RoomList[key].Event.Start)

	for i := dayIdx * 24; i < (dayIdx+1)*24 && i < len(schedule); i++ {
		timeOutput := startTime.Add(time.Hour * time.Duration(i)).Format(discord.TimeOutputFormat)

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
