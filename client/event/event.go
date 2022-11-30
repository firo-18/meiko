package event

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/api"
	"github.com/firo-18/meiko/schema"
)

var (
	List       = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	EventList  = make(map[string]schema.Event)
	FillerList = make(map[string]schema.Filler)
	RoomList   = make(map[string]*schema.Room)
)

func init() {

	schema.DeserializeRooms(&RoomList)
	schema.DeserializeFillers(&FillerList)

	go fetchEvents()
}

func fetchEvents() {
	url := "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-en-diff/main/events.json"
	err := api.GetDecodeJSON(url, EventList)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(EventList)
}
