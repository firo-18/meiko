package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/schema"
)

var (
	ISO8601          = "2006-01-02T03:04:05-0700"
	TimeOutputFormat = "Jan 02 15:04"
	EmbedTimestamp   = time.Now().Format(ISO8601)
	EmbedColor       = 15548997
	EmbedFooter      = func(s *discordgo.Session) *discordgo.MessageEmbedFooter {
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

func RoomInfoFields(room *schema.Room) []*discordgo.MessageEmbedField {
	return []*discordgo.MessageEmbedField{
		{
			Name:  "Server",
			Value: StyleFieldValues(room.Server),
		},
		{
			Name:  "Event",
			Value: StyleFieldValues(room.Event.Name),
		},
		{
			Name:  "Created By",
			Value: StyleFieldValues(room.Owner.Username),
		},
		{
			Name:  "Runner",
			Value: StyleFieldValues(room.Runner),
		},
	}
}
