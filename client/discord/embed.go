package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/ghost"
)

var (
	ISO8601        = "2006-01-02T03:04:05-0700"
	EmbedTimestamp = time.Now().Format(ISO8601)
	EmbedColor     = 15548997
	EmbedFooter    = func(s *discordgo.Session) *discordgo.MessageEmbedFooter {
		return &discordgo.MessageEmbedFooter{
			Text:    s.State.User.Username,
			IconURL: s.State.User.AvatarURL(""),
		}
	}
)

func NewEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     EmbedColor,
		Timestamp: EmbedTimestamp,
		Footer:    EmbedFooter(s),
	}
}

func NewListEmbed(s *discordgo.Session, list map[string]*ghost.Ghost) []*discordgo.MessageEmbed {
	return []*discordgo.MessageEmbed{
		{
			Title:     "List",
			Color:     EmbedColor,
			Timestamp: EmbedTimestamp,
			Footer:    EmbedFooter(s),
			Fields:    listFields(list),
		},
	}
}

func listFields(list map[string]*ghost.Ghost) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}

	for _, v := range list {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   v.User.Username,
			Value:  StyleFieldValues(v.SkillValue),
			Inline: true,
		})
	}
	return fields
}
