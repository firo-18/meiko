package schema

import (
	"encoding/gob"
	"encoding/json"
	"log"
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
	Schedule    [][]*Filler    `json:"schedule"`
	Owner       discordgo.User `json:"owner"`
	Manager     discordgo.User `json:"manager"`
	CreateAt    time.Time      `json:"createAt"`
}

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

// Backup writes room data to local json file for back up.
func (r *Room) Backup() error {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}

	filename := PathRoomDB + r.Key + " - " + r.CreateAt.String() + ".json"
	err = os.WriteFile(filename, data, 0640)
	return err
}

// SerializeRooms encodes all rooms data into local a local gob file when client end, intentional or not.
func SerializeRooms(rooms map[string]*Room) {
	filename := PathDB + "rooms.gob"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := gob.NewEncoder(f)
	err = data.Encode(rooms)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success: Rooms data has been serialized.")
}

// DeserializeRooms decodes rooms data when client starts from local gob file to memory.
func DeserializeRooms(rooms *map[string]*Room) {
	filename := PathDB + "rooms.gob"
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if stat.Size() == 0 {
		return
	}

	data := gob.NewDecoder(f)
	err = data.Decode(rooms)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success: Rooms data has been deserialized.")
}
