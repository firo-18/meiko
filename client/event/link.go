package event

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
	"github.com/firo-18/meiko/schema"
)

func init() {
	List["link"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			user := i.Member.User
			var lead, sum int64
			var offset string

			for i, v := range data.Options {
				switch v.Name {
				case "lead":
					lead = data.Options[i].IntValue()
				case "sum":
					sum = data.Options[i].IntValue()
				case "utc-offset":
					offset = data.Options[i].StringValue()
				}
			}

			offsetNum, err := strconv.Atoi(offset)
			if err != nil {
				discord.EmbedError(s, i, discord.EmbedErrInvalidOffset)
				return
			}

			if offsetNum < -12 || offsetNum > 14 {
				discord.EmbedError(s, i, discord.EmbedErrInvalidOffset)
				return
			}

			// Calculate skill multiplier from ISV.
			skillValue := (float64(sum-lead) * 0.002) + float64(lead)/100 + 1

			// Add or update filler.
			FillerList[user.ID] = schema.NewFiller(user, skillValue, offsetNum)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Success",
							Description: "You have join the filler list. If your ISV changes, re-run this command to update.",
							Color:       discord.EmbedColor,
							Timestamp:   discord.EmbedTimestamp,
							Footer:      discord.EmbedFooter(s),
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "ID",
									Value: discord.StyleFieldValues(user.ID),
								},
								{
									Name:  "UTC Offset",
									Value: discord.StyleFieldValues(offset),
								},
								{
									Name:  "ISV",
									Value: discord.StyleFieldValues(lead, "/", sum),
								},
								{
									Name:  "Skill Multiplier Value",
									Value: discord.StyleFieldValues(skillValue),
								},
							},
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			// Log link activities.
			log.Printf("%v has linked to %v: %v %v", user.Username, s.State.User.String(), skillValue, offset)

		// Autocomplete UTC offset.
		case discordgo.InteractionApplicationCommandAutocomplete:
			offsetAutocomplete(s, i)
		}
	}
}

func offsetAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	var choice string

	for j, v := range data.Options {
		if v.Name == "utc-offset" {
			choice = data.Options[j].StringValue()
		}
	}

	offsets := []string{}

	for i := -12; i <= 12; i++ {
		currentUTC := time.Now().UTC()
		offset := currentUTC.Add(time.Hour * time.Duration(i))
		timeString := offset.Format(discord.TimeOutputFormat)
		offsets = append(offsets, fmt.Sprint("UTC ", i, " - ", timeString))
	}

	for j, offset := range offsets {
		if ok, _ := regexp.MatchString("(?i)"+choice, offset); ok {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  offset,
				Value: fmt.Sprint(j - 12),
			})
		}
	}

	// Max number of choice is 25.
	if len(choices) > 25 {
		choices = choices[:25]
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
