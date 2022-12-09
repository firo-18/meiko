package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/schema"
)

func DayScheduleComponents(room *schema.Room, filler *schema.Filler, cmd string) []discordgo.MessageComponent {
	options := []discordgo.SelectMenuOption{}

	dSum := room.EventLength/24 + 1
	for d := 0; d < dSum; d++ {

		if cmd != "view" {
			// Find the last hour of the event day time.
			dayLastHour := (d+1)*24 - 1
			eventDayLastHour := time.UnixMilli(room.Event.Start).Add(time.Hour * time.Duration(dayLastHour))

			if time.Now().After(eventDayLastHour) {
				continue
			}
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:       fmt.Sprint("Day ", d+1),
			Value:       fmt.Sprint(room.Key, "_", d),
			Description: "Start from " + time.UnixMilli(room.Event.Start).Add(time.Hour*24*time.Duration(d)).Add(time.Hour*time.Duration(filler.Offset)).UTC().Format(TimeOutputFormat) + " offset time.",
		})
	}

	menu := discordgo.SelectMenu{
		CustomID:    "menu-day",
		Placeholder: "Select a day.",
		Options:     options,
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

func HourScheduleComponents(room *schema.Room, filler *schema.Filler, d int, cmd string) []discordgo.MessageComponent {
	options := []discordgo.SelectMenuOption{}

	startTime := time.UnixMilli(room.Event.Start)
	maxOption := 1

	switch cmd {
	case "schedule":
		options = append(options, discordgo.SelectMenuOption{
			Label:       "Deselect All",
			Value:       fmt.Sprint(room.Key, "_", d, "_", d*24, "_", "deselect"),
			Description: "Useful for when deselecting all hours.",
			Default:     false,
		})

		for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
			eventTime := startTime.Add(time.Hour * time.Duration(h))
			if eventTime.After(time.Now()) {
				maxOption++
				_, ok := schema.HasShift(room.Schedule[h], filler.User.ID)

				options = append(options, discordgo.SelectMenuOption{
					Label:       eventTime.Add(time.Hour * time.Duration(filler.Offset)).UTC().Format(TimeOutputFormat),
					Value:       fmt.Sprint(room.Key, "_", d, "_", h),
					Description: fmt.Sprint("Event Hour: ", h),
					Default:     ok,
				})
			}
		}

	case "manage":
		options = append(options, discordgo.SelectMenuOption{
			Label:       "Back",
			Value:       fmt.Sprint(room.Key, "_", d, "_", d*24, "_", "back"),
			Description: "Go back to the previous options.",
			Default:     false,
		})

		for h := d * 24; h < (d+1)*24 && h < len(room.Schedule); h++ {
			eventTime := startTime.Add(time.Hour * time.Duration(h))
			if eventTime.After(time.Now()) {
				if len(room.Schedule[h]) > 0 {
					options = append(options, discordgo.SelectMenuOption{
						Label:       eventTime.Add(time.Hour * time.Duration(filler.Offset)).UTC().Format(TimeOutputFormat),
						Value:       fmt.Sprint(room.Key, "_", d, "_", h),
						Description: fmt.Sprint("Event Hour: ", h),
					})
				}
			}
		}

	}

	menu := discordgo.SelectMenu{
		CustomID:    "menu-hour",
		Placeholder: "Select a shift.",
		Options:     options,
		MaxValues:   maxOption,
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

func FillerScheduleComponents(fillerList map[string]*schema.Filler, room *schema.Room, h int) []discordgo.MessageComponent {
	options := []discordgo.SelectMenuOption{
		{
			Label:       "Default",
			Value:       fmt.Sprint(room.Key, "_", h, "_", "default"),
			Description: "Useful for when removing all fillers.",
			Default:     true,
		},
	}

	for _, fillerID := range room.Schedule[h] {
		filler := fillerList[fillerID]
		options = append(options, discordgo.SelectMenuOption{
			Label:       filler.User.Username,
			Description: fmt.Sprintf("%v (%.2f)", filler.ISV, filler.SkillValue),
			Value:       fmt.Sprint(room.Key, "_", h, "_", filler.User.ID),
			Default:     true,
		})
	}

	menu := discordgo.SelectMenu{
		CustomID:    "menu-fillers",
		Placeholder: "You can only remove filler(s). Irreversible.",
		MaxValues:   len(room.Schedule[h]) + 1,
		Options:     options,
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
