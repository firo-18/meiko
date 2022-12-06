package schema

import (
	"encoding/json"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Filler defines the type of data stores for each Filler.
type Filler struct {
	User         discordgo.User `json:"user"`
	ISV          string         `json:"isv"`
	SkillValue   float64        `json:"skillValue"`
	Offset       int            `json:"offset"`
	LastModified time.Time      `json:"lastModified"`
	CreatedAt    time.Time      `json:"createdAt"`
}

// New takes constructs a new Ghost with initial values from a discordgo.User input and return its address.
func NewFiller(user *discordgo.User, isv string, skill float64, offset int) *Filler {
	return &Filler{
		User:         *user,
		ISV:          isv,
		SkillValue:   skill,
		Offset:       offset,
		LastModified: time.Now(),
		CreatedAt:    time.Now(),
	}
}

// RestoreFiller restores filler data from local file and returns the Filler's address.
func RestoreFiller(filename string) (*Filler, error) {
	file, err := os.ReadFile(PathFillerDB + filename)
	if err != nil {
		return nil, err
	}
	var f Filler
	err = json.Unmarshal(file, &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
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
