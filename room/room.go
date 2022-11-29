package room

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/event"
	"github.com/firo-18/meiko/ghost"
)

const (
	UnixMilliPerHour   = 1_000 * 3_600
	EventEndTimeOffset = 1_000
)

// Room lists all the scheduling and user data for tiering.
type Room struct {
	Key         string           `json:"key"`
	Name        string           `json:"name"`
	Server      string           `json:"server"`
	Runner      string           `json:"runner"`
	Event       event.Event      `json:"event"`
	EventLength int              `json:"length"`
	Fillers     []ghost.Ghost    `json:"fillers"`
	Schedule    [][]*ghost.Ghost `json:"schedule"`
	Creator     discordgo.User   `json:"creator"`
	CreateAt    time.Time        `json:"createAt"`
}

func New(e event.Event, creator discordgo.User, name, guildID string) *Room {
	length := int(e.End-e.Start+EventEndTimeOffset) / UnixMilliPerHour
	return &Room{
		Key:         guildID + " - " + name,
		Name:        name,
		Server:      guildID,
		Event:       e,
		EventLength: length,
		Creator:     creator,
		Runner:      creator.Username,
		Schedule:    make([][]*ghost.Ghost, length),
		CreateAt:    time.Now(),
	}
}
