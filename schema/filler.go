package schema

import (
	"encoding/gob"
	"encoding/json"
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
	LastModified time.Time      `json:"lastModified"`
}

// New takes constructs a new Ghost with initial values from a discordgo.User input and return its address.
func NewFiller(user *discordgo.User, skill float64, offset int) *Filler {
	return &Filler{
		User:         *user,
		SkillValue:   skill,
		Offset:       offset,
		LastModified: time.Now(),
	}
}

// Backup writes filler data to local json file for back up.
func (f *Filler) Backup() error {
	data, err := json.MarshalIndent(f, "", "\t")
	if err != nil {
		return err
	}

	filename := PathFillerDB + f.User.String() + ".json"
	err = os.WriteFile(filename, data, 0640)
	return err
}

// SerializeFillers encodes all fillers data into local a local gob file when client end, intentional or not.
func SerializeFillers(fillers map[string]*Filler) {
	filename := PathDB + "fillers.gob"
	f, err := os.Create(filename)
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

// DeserializeFillers decodes fillers data when client starts from local gob file to memory.
func DeserializeFillers(fillers *map[string]*Filler) {
	filename := PathDB + "fillers.gob"

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
	err = data.Decode(fillers)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success: Fillers data has been deserialized.")
}
