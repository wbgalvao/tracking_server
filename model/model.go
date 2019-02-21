package model

import (
	"time"

	"github.com/gocql/gocql"
)

// Event holds information about an event on the internet.
// 	- Identification is a field designed to capture the user identification
// 	- Target is a field designed to capture the target web page  address
// 	- Event will define what type of event the Event captured
// 	- Valid is designed to validate Event instances based on the Event field
// 	- Timestamp will hold the time of the event's occurrence
type Event struct {
	ID          gocql.UUID `json:"id,omitempty"`
	Username    string     `json:"username,omitempty"`
	Target      string     `json:"target,omitempty"`
	Description string     `json:"description,omitempty"`
	Valid       bool       `json:"valid,omitempty"`
	Timestamp   time.Time  `json:"timestamp,omitempty"`
}

// ValidEvents defines the events specification which will be considered
// valid for the application.
var ValidEvents = [2]string{"CLICK", "VIEW"}

// Validate uses the ValidEvents array to check if the built Event has a valid
// Description field.
func (e *Event) Validate() {
	for _, event := range ValidEvents {
		if e.Description == event {
			e.Valid = true
			e.Timestamp = time.Now()
			e.ID = gocql.TimeUUID()
			return
		}
	}
	e.Valid = false
}
