package ghost

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Ghost defines the type of data stores for each User.
type Ghost struct {
	User         discordgo.User `json:"user"`
	SkillValue   float64        `json:"skillValue"`
	CreatedAt    time.Time      `json:"createdAt"`
	LastModified time.Time      `json:"lastModified"`
}

// New takes constructs a new Ghost with initial values from a discordgo.User input and return its address.
func New(user discordgo.User, skill float64) Ghost {
	return Ghost{
		User:         user,
		SkillValue:   skill,
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}
}
