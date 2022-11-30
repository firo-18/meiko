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
	Key         string             `json:"key"`
	Name        string             `json:"name"`
	Server      string             `json:"server"`
	Runner      string             `json:"runner"`
	Event       Event              `json:"event"`
	EventLength int                `json:"length"`
	FillerList  map[string]*Filler `json:"fillerList"`
	Schedule    [][]*Filler        `json:"schedule"`
	Owner       discordgo.User     `json:"owner"`
	CreateAt    time.Time          `json:"createAt"`
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
		Runner:      owner.Username,
		FillerList:  make(map[string]*Filler),
		Schedule:    make([][]*Filler, length),
		CreateAt:    time.Now(),
	}
}

// Backup encodes room data to a local gob file. Use for persistently update room data.
func (r *Room) Backup() error {
	path := "./db/rooms/"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(path+r.Key+".gob", os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(r); err != nil {
		return err
	}
	log.Print("Successfully encode room to local file.")
	return nil
}

func (r *Room) Archive() error {
	path := "./db/rooms/"

	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(path+r.Key+".json", data, 0640)
	return err
}

func SerializeRooms(rooms map[string]*Room) {
	f, err := os.Create("rooms.gob")
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

func DeserializeRooms(rooms *map[string]*Room) {
	f, err := os.OpenFile("rooms.gob", os.O_RDONLY|os.O_CREATE, 0640)
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
