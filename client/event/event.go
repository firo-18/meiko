package event

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/api"
	"github.com/firo-18/meiko/db"
	"github.com/firo-18/meiko/event"
	"github.com/firo-18/meiko/room"
)

var (
	List      = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	EventList = make([]event.Event, 10)
	RoomList  = make(map[string]*room.Room)
)

func init() {
	db.DeserializeRooms(&RoomList)
	go fetchEvents()
	log.Println(RoomList)
	// log.Println(len(RoomList["TEST"].Fillers))
	// log.Println(RoomList["TEST"].Fillers[0].SkillValue)
}

func fetchEvents() {
	url := "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-en-diff/main/events.json"
	err := api.GetDecodeJSON(url, &EventList)
	if err != nil {
		log.Fatal(err)
	}
}
