package event

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/discord"
)

func init() {
	List["stats"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			guildList := []string{}
			for _, g := range s.State.Guilds {
				guildList = append(guildList, g.Name)
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:     fmt.Sprint(s.State.User.String(), " - Stats"),
							Color:     discord.EmbedColor,
							Timestamp: discord.EmbedTimestamp,
							Footer:    discord.EmbedFooter(s),
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "In Guilds",
									Value: discord.StyleFieldValues(len(s.State.Guilds)),
								},
								{
									Name:  "Guild List",
									Value: discord.StyleFieldValues(strings.Join(guildList, ", ")),
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

		// Autocomplete UTC offset.
		case discordgo.InteractionApplicationCommandAutocomplete:
			offsetAutocomplete(s, i)
		}
	}
}
