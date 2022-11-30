package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/firo-18/meiko/schema"
)

// ReadAPI takes a URL string argument and a container to store decoded JSON data, and return an error if response code is not 200.
func GetDecodeJSON(url string, eventList map[string]schema.Event) error {
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
		var event schema.Event
		err := dec.Decode(&event)
		if err != nil {
			panic(err)
		}
		fmt.Printf("OK: event=%#v\n", event)

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
