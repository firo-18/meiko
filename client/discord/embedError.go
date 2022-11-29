package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	EmbedErrorRoom404            = 1
	EmbedErrorRoomNotJoined      = 2
	EmbedErrorRoomEnded          = 3
	EmbedErrorRoomNameDuplicated = 4
	EmbedErrorInvalidInteraction = 5
)

func EmbedError(s *discordgo.Session, i *discordgo.InteractionCreate, code int) {
	embed := NewEmbed(s)
	embed.Title = "Error"

	switch code {
	case 1:
		embed.Description = "Room not exist. Select a room from the list."
	case 2:
		embed.Description = "You have not join the room yet. Use /join command before scheduling your hours."
	case 3:
		embed.Description = "Event has ended. Room data will be archived shortly."
	case 4:
		embed.Description = "Room name already exists. Choose a different name."
	case 5:
		embed.Description = "This interaction is intended for the original user only."
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
