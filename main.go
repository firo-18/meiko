package main

import (
	"log"
	"os"

	"github.com/firo-18/meiko/client"
	"github.com/firo-18/meiko/schema"
)

func init() {
	// Mkdir all neccessary path
	if err := os.MkdirAll(schema.PathRoomArchive, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(schema.PathFillerDB, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// command.DeployTest()
	// command.DeployProduction()
	client.Open()
}
