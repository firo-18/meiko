package db

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/firo-18/meiko/room"
)

func SerializeRooms(rooms map[string]*room.Room) {
	f, err := os.Create("rooms.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := gob.NewEncoder(f)
	err = data.Encode(rooms)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success: Rooms data has been serialized.")
}

func DeserializeRooms(rooms *map[string]*room.Room) {
	f, err := os.OpenFile("rooms.gob", os.O_RDONLY|os.O_CREATE, 0640)
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
	err = data.Decode(rooms)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success: Rooms data has been deserialized.")
}
