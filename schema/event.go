package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Event defines Project Sekai event data.
type Event struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"eventType"`
	Start int64  `json:"startAt"`
	End   int64  `json:"aggregateAt"`
}

// GetEventList fetches event data from JSON API, decodes only upcoming events into EventList, and returns any error.
func GetEventList(url string, eventList map[string]Event) error {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("http-get:", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	dec := json.NewDecoder(res.Body)

	// Read the open bracket.
	t, err := dec.Token()
	if err != nil {
		panic(err)
	}
	fmt.Printf("OK: %T: %v\n", t, t)

	// While the array contains values.
	for dec.More() {
		// Decode an array value.
		var event Event
		err := dec.Decode(&event)
		if err != nil {
			panic(err)
		}
		if time.Now().Before(time.UnixMilli(event.End)) {
			eventList[event.Name] = event
			log.Printf("Loaded event '%v' into event list.", event.Name)
		}
	}

	// Read the closing bracket.
	t, err = dec.Token()
	if err != nil {
		panic(err)
	}
	fmt.Printf("OK: %T: %v\n", t, t)

	return nil
}
