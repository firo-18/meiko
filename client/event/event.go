package event

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/db"
	"github.com/firo-18/meiko/schema"
)

var (
	List       = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	EventList  = make(map[string]schema.Event)
	FillerList = make(map[string]*schema.Filler)
	RoomList   = make(map[string]map[string]*schema.Room)
)

func init() {
	// Mkdir all neccessary path
	if err := os.MkdirAll(schema.PathRoomDB, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(schema.PathFillerDB, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	var err error
	RoomList, err = db.FetchRoomList()
	if err != nil {
		LogError(err, "fetchRoomList")
	}

	FillerList, err = db.FetchFillers()
	if err != nil {
		LogError(err, "fetchFillers")
	}

	go fetchEvents()
}

func fetchEvents() {
	url := "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-en-diff/main/events.json"
	err := schema.GetEventList(url, EventList)
	if err != nil {
		log.Fatal(err)
	}
}

func LogError(err error, cmd string) {
	log.Printf("Origin: %v. Error: %v", cmd, err)
}
