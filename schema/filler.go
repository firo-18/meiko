package schema

import (
	"encoding/gob"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Filler defines the type of data stores for each Filler.
type Filler struct {
	User         discordgo.User `json:"user"`
	SkillValue   float64        `json:"skillValue"`
	Offset       int            `json:"offset"`
	CreatedAt    time.Time      `json:"createdAt"`
	LastModified time.Time      `json:"lastModified"`
}

// New takes constructs a new Ghost with initial values from a discordgo.User input and return its address.
func NewFiller(user *discordgo.User, skill float64, offset int) Filler {
	return Filler{
		User:         *user,
		SkillValue:   skill,
		Offset:       offset,
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}
}

func SerializeFillers(fillers map[string]Filler) {
	f, err := os.Create("fillers.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := gob.NewEncoder(f)
	err = data.Encode(fillers)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success: Fillers data has been serialized.")
}

func DeserializeFillers(fillers *map[string]Filler) {
	f, err := os.OpenFile("fillers.gob", os.O_RDONLY|os.O_CREATE, 0640)
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
	err = data.Decode(fillers)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success: Fillers data has been deserialized.")
}
