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
	Guild       string         `json:"guild"`
	Event       Event          `json:"event"`
	EventLength int            `json:"length"`
	Owner       discordgo.User `json:"owner"`
	Manager     discordgo.User `json:"manager"`
	CreateAt    time.Time      `json:"createAt"`
	Schedule    [][]string     `json:"schedule"`
}

// NewRoom creates a new Room based on arguments, and return the Room's address.
func NewRoom(guildID, name string, event Event, owner *discordgo.User) *Room {
	length := int(event.End-event.Start+EventEndTimeOffset) / UnixMilliPerHour
	return &Room{
		Key:         guildID + "_" + name,
		Name:        name,
		Guild:       guildID,
		Event:       event,
		EventLength: length,
		Owner:       *owner,
		Manager:     *owner,
		Schedule:    make([][]string, length),
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

	err = os.MkdirAll(PathRoomDB+r.Guild+"/archive/", os.ModePerm)
	if err != nil {
		return err
	}

	filename := PathRoomDB + r.Guild + "/" + r.Name + ".json"
	err = os.WriteFile(filename, data, 0640)
	return err
}

// Archive writes room data to local json file for archive-purpose.
func (r *Room) Archive() error {
	r.Backup()

	oldFilename := PathRoomDB + r.Guild + "/" + r.Name + ".json"
	newFilename := PathRoomDB + r.Guild + "/archive/" + r.Name + " - " + r.CreateAt.String() + ".json"
	err := os.Rename(oldFilename, newFilename)

	return err
}

// GetFillers loop through Room's schedule and return all unique fillers' discord User.
func (r *Room) GetFillers(fillerList map[string]*Filler) map[string]*Filler {
	list := make(map[string]*Filler)

	for _, h := range r.Schedule {
		for _, id := range h {
			if filler, ok := fillerList[id]; ok {
				list[id] = filler
			}
		}
	}
	return list
}

func HasShift(fillers []string, userID string) (int, bool) {
	for i, filler := range fillers {
		if filler == userID {
			return i, true
		}
	}
	return 0, false
}
