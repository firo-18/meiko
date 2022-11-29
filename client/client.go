package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/component"
	"github.com/firo-18/meiko/client/config"
	"github.com/firo-18/meiko/client/event"
	"github.com/firo-18/meiko/db"
)

var (
	s *discordgo.Session
)

func init() {
	cfx := config.Load("config.json")

	var err error
	s, err = discordgo.New("Bot " + cfx.Token)
	if err != nil {
		log.Fatalln(`Error setting up discord client:`, err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("%v#%v is up and running...", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			if h, ok := event.List[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := component.List[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
}

func Open() {
	err := s.Open()
	if err != nil {
		log.Fatalln("Cannot open the session:", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	log.Println("Press Ctrl+C to exit")

	<-stop

	close()

}

// close serialzes internal data to files before closing websocket and exit the program.
func close() {
	db.SerializeRooms(event.RoomList)

	log.Println("Gracefully shutting down.")
}
