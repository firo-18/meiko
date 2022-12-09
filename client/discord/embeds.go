package discord

import (
	"fmt"
	"strings"
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
			Name:  "Event",
			Value: FieldStyle(room.Event.Name),
		},
		{
			Name:  "Event Type",
			Value: FieldStyle(room.Event.Type),
		},
		{
			Name:  "Event Length",
			Value: FieldStyle(room.EventLength),
		},
		{
			Name:  "Owner",
			Value: FieldStyle(room.Owner.Username),
		},
		{
			Name:  "Manager",
			Value: FieldStyle(room.Manager.Username),
		},
	}
}

func DayScheduleEmbeds(s *discordgo.Session, fillerList map[string]*schema.Filler, room *schema.Room, filler *schema.Filler, d int, cmd string) []*discordgo.MessageEmbed {
	eventStart := time.UnixMilli(room.Event.Start)

	fields := []*discordgo.MessageEmbedField{}

	for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
		eventTime := eventStart.Add(time.Hour * time.Duration(h))

		if cmd != "view" {
			if eventTime.Before(time.Now()) {
				continue
			}
		}

		fillers := make([]string, len(room.Schedule[h]))
		for j, v := range room.Schedule[h] {
			shiftFiller := fillerList[v]
			if shiftFiller.User.ID == filler.User.ID {
				fillers[j] = fmt.Sprintf("__%v__", shiftFiller.User.Username)
			} else {
				fillers[j] = shiftFiller.User.Username
			}
		}
		value := strings.Join(fillers, ", ")
		if value == "" {
			value = "-"
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Hour %v: <t:%v:t>", h, eventTime.Unix()),
			Value:  FieldStyle(value),
			Inline: true,
		})
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Title:     fmt.Sprint("[Day ", d+1, "] Room - ", room.Name),
			Color:     EmbedColor,
			Timestamp: EmbedTimestamp,
			Footer:    EmbedFooter(s),
			Fields:    fields,
		},
	}

	return embeds
}

func FillerScheduleEmbeds(s *discordgo.Session, fillerList map[string]*schema.Filler, room *schema.Room, d, h int) []*discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{}

	for _, fillerID := range room.Schedule[h] {
		filler := fillerList[fillerID]
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  filler.User.Username,
			Value: FieldStyle(fmt.Sprintf("%v (%.2f)", filler.ISV, filler.SkillValue)),
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("Day %v - Hour %v", d+1, h),
		Color:     EmbedColor,
		Timestamp: EmbedTimestamp,
		Footer:    EmbedFooter(s),
		Fields:    fields,
	}

	embeds := []*discordgo.MessageEmbed{embed}

	return embeds
}
