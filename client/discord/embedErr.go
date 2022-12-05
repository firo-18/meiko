package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	EmbedErrInvalidFiller      = 1
	EmbedErrInvalidInteraction = 2
	EmbedErrInvalidEvent       = 3
	EmbedErrRoom404            = 4
	EmbedErrRoomNameDuplicated = 5
	EmbedErrRoomEnded          = 6
	EmbedErrInvalidOffset      = 7
	EmbedErrInvalidOwner       = 8
)

func EmbedError(s *discordgo.Session, i *discordgo.InteractionCreate, code int) {
	embed := NewEmbed(s)
	embed.Title = "Error"

	switch code {
	case 1:
		embed.Description = "You have not linked your ISV yet. Run /link first."
	case 2:
		embed.Description = "This interaction is intended for the original user only."
	case 3:
		embed.Description = "Invalid event name. Select an event from the available pool."
	case 4:
		embed.Description = "Room not exist. Select a room from the list."
	case 5:
		embed.Description = "Room name already exists. Choose a different name."
	case 6:
		embed.Description = "Event has ended. Room data will be archived shortly."
	case 7:
		embed.Description = "Invalid offset. Select from the option, or enter an integer between -12 and 12, inclusive."
	case 8:
		embed.Description = "You are not the owner/manager of this room. Only owner/manager can alter the room."
	}

	embeds := []*discordgo.MessageEmbed{embed}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
