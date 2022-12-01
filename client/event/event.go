package event

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/api"
	"github.com/firo-18/meiko/schema"
)

var (
	List       = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	EventList  = make(map[string]schema.Event)
	FillerList = make(map[string]*schema.Filler)
	RoomList   = make(map[string]*schema.Room)
)

func init() {
	// Mkdir all neccessary path
	if err := os.MkdirAll(schema.PathRoomArchive, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(schema.PathFillerDB, os.ModePerm); err != nil {
		log.Fatal(err)
	}

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

func errorRestart(err error) {
	log.Println("Client restarting due to error encountered:", err)
	schema.SerializeRooms(RoomList)
	schema.SerializeFillers(FillerList)

	stop := make(chan os.Signal, 1)
	<-stop

}
