package event

import (
	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/ghost"
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
