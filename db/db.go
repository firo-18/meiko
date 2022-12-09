package db

import (
	"os"
	"strings"
	"time"

	"github.com/firo-18/meiko/schema"
)

func FetchRoomList() (map[string]map[string]*schema.Room, error) {
	roomList := make(map[string]map[string]*schema.Room)

	files, err := os.ReadDir(schema.PathRoomDB)
	if err != nil {
		return nil, err
	}

	for _, fs := range files {
		if fs.IsDir() {
			if list, err := FetchRoomsInGuild(fs.Name()); err != nil {
				return nil, err
			} else {
				roomList[fs.Name()] = list
			}
		}
	}

	return roomList, nil
}

func FetchRoomsInGuild(server string) (map[string]*schema.Room, error) {
	roomList := make(map[string]*schema.Room)

	path := schema.PathRoomDB + server + "/"
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			if r, err := schema.RestoreRoom(server, file.Name()); err != nil {
				return nil, err
			} else {
				if time.Now().After(time.UnixMilli(r.Event.End)) {
					err := r.Archive()
					if err != nil {
						return nil, err
					}
				} else {
					roomList[r.Name] = r
				}
			}
		}
	}

	return roomList, nil
}

func FetchFillers() (map[string]*schema.Filler, error) {
	fillerList := make(map[string]*schema.Filler)

	files, err := os.ReadDir(schema.PathFillerDB)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			if r, err := schema.RestoreFiller(file.Name()); err != nil {
				return nil, err
			} else {
				fillerList[r.User.ID] = r
			}
		}
	}

	return fillerList, nil
}
