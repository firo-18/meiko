package event

// Event defines Project Sekai event data.
type Event struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"eventType"`
	Start int64  `json:"startAt"`
	End   int64  `json:"aggregateAt"`
}

func New() *Event {
	return &Event{}
}
