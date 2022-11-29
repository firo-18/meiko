package command

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/firo-18/meiko/client/config"
)

var (
	List                         = []*discordgo.ApplicationCommand{}
	permissionManageServer int64 = discordgo.PermissionManageServer
)

// DeployProd deploys all commands into production discord bot.
func DeployProduction() {
	config := config.Load("config.json")

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalln(`Error setting up discord client:`, err)
	}

	registeredCommands, err := s.ApplicationCommandBulkOverwrite(config.ClientID, "", List)
	if err != nil {
		log.Fatalln("Cannot create commands:", err)
	}
	log.Printf("Deployed %v slash production commands successfully to all servers.", len(registeredCommands))
}

// DeployProd deploys all commands into production discord bot.
func DeployTest() {
	config := config.Load("config.json")

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalln(`Error setting up discord client:`, err)
	}

	registeredCommands, err := s.ApplicationCommandBulkOverwrite(config.ClientID, config.GuildID, List)
	if err != nil {
		log.Fatalln("Cannot create commands:", err)
	}
	log.Printf("Deployed %v slash test commands successfully to test server.", len(registeredCommands))
}
