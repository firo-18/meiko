package schema

import (
	"encoding/json"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	UnixMilliPerHour   = 1_000 * 3_600
	EventEndTimeOffset = 1_000
)

// Room lists all the scheduling and user data for tiering.
type Room struct {
	Key         string         `json:"key"`
	Name        string         `json:"name"`
	Server      string         `json:"server"`
	Event       Event          `json:"event"`
	EventLength int            `json:"length"`
	Owner       discordgo.User `json:"owner"`
	Manager     discordgo.User `json:"manager"`
	CreateAt    time.Time      `json:"createAt"`
	Schedule    [][]*Filler    `json:"schedule"`
}

// NewRoom creates a new Room based on arguments, and return the Room's address.
func NewRoom(guildID, name string, event Event, owner *discordgo.User) *Room {
	length := int(event.End-event.Start+EventEndTimeOffset) / UnixMilliPerHour
	return &Room{
		Key:         guildID + " - " + name,
		Name:        name,
		Server:      guildID,
		Event:       event,
		EventLength: length,
		Owner:       *owner,
		Manager:     *owner,
		Schedule:    make([][]*Filler, length),
		CreateAt:    time.Now(),
	}
}

// RestoreRoom restores room data from local file and returns the Room's address.
func RestoreRoom(key, filename string) (*Room, error) {
	file, err := os.ReadFile(PathRoomDB + key + "/" + filename)
	if err != nil {
		return nil, err
	}
	var r Room
	err = json.Unmarshal(file, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Backup writes room data to local json file for back up.
func (r *Room) Backup() error {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}

	err = os.MkdirAll(PathRoomDB+r.Server+"/archive/", os.ModePerm)
	if err != nil {
		return err
	}

	filename := PathRoomDB + r.Server + "/" + r.Name + ".json"
	err = os.WriteFile(filename, data, 0640)
	return err
}

// Archive writes room data to local json file for archive-purpose.
func (r *Room) Archive() error {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}

	filename := PathRoomDB + r.Server + "/archive/" + r.Name + " - " + r.CreateAt.String() + ".json"
	err = os.WriteFile(filename, data, 0640)
	return err
}
