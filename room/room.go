package room

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/event"
	"github.com/firo-18/meiko/ghost"
)

const (
	UnixPerHour        = 1_000 * 3_600
	EventEndTimeOffset = 1_000
)

// Room lists all the scheduling and user data for tiering.
type Room struct {
	Event    event.Event
	Creator  discordgo.User
	Runner   string
	Fillers  []*ghost.Ghost
	Schedule [][]*ghost.Ghost
	CreateAt time.Time
}

func New(e event.Event, creator discordgo.User) *Room {
	return &Room{
		Event:    e,
		Creator:  creator,
		Runner:   creator.Username,
		Schedule: make([][]*ghost.Ghost, (e.End-e.Start+EventEndTimeOffset)/UnixPerHour),
		CreateAt: time.Now(),
	}
}
