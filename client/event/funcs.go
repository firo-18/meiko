package event

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/ghost"
	"github.com/firo-18/meiko/room"
)

var (
	MaxOption        = 25
	OptionDateFormat = "Jan 02 15:04"
)

// findFiller loops through the fillers slice to find a filler. If exists, it returns the index and true, otherwise, it returns 0 and false.
func FindFiller(fillers []ghost.Ghost, filler discordgo.User) (int, bool) {
	for i, v := range fillers {
		if v.User == filler {
			return i, true
		}
	}
	return 0, false
}

func HasShift(userID, roomName string, hour int) bool {
	for _, filler := range RoomList[roomName].Schedule[hour] {
		if filler.User.ID == userID {
			return true
		}
	}
	return false
}

func ScheduleDayComponent(room *room.Room, days int) []discordgo.MessageComponent {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "day-select",
					Placeholder: "Select a day to schedule.",
					Options:     DaySelectMenu(room, days),
				},
			},
		},
	}

	return components
}

func DaySelectMenu(room *room.Room, days int) []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{}

	for i := 0; i < days; i++ {
		options = append(options, discordgo.SelectMenuOption{
			Label:       fmt.Sprint("Day ", i+1),
			Value:       fmt.Sprint(room.Key, " - ", i+1),
			Description: time.UnixMilli(room.Event.Start).Add(time.Hour * 24 * time.Duration(i)).Local().Format(OptionDateFormat),
		})
	}
	return options
}

func ScheduleMenu(arr []time.Time, start int64, roomName, componentID, userID string) discordgo.SelectMenu {
	options := make([]discordgo.SelectMenuOption, len(arr))

	for i, v := range arr {
		hour := int(v.Sub(time.UnixMilli(start)).Hours())
		options[i] = discordgo.SelectMenuOption{
			Label:       v.Format("Jan 02 15:04"),
			Value:       fmt.Sprintf("%v_%v", roomName, hour),
			Description: "Local time.",
			Default:     HasShift(userID, roomName, hour),
		}
	}

	menu := discordgo.SelectMenu{
		CustomID:    componentID,
		Placeholder: "Select your available time slot.",
		MaxValues:   len(arr),
		Options:     options,
	}

	return menu
}
