package model

import (
	"testing"
	"time"

	"github.com/gocql/gocql"
)

func TestValidateEvent(t *testing.T) {
	e := Event{
		ID:          gocql.TimeUUID(),
		Username:    "John Doe",
		Target:      "https://www.google.com/",
		Description: "CLICK",
		Timestamp:   time.Now(),
	}
	e.Validate()
	if e.Valid == false {
		t.Errorf("error validating event tracker. expected: %t, got: %t\n",
			true, e.Valid)
	}
}

func TestInvalidEvent(t *testing.T) {
	e := Event{
		ID:          gocql.TimeUUID(),
		Username:    "John Doe",
		Target:      "https://www.google.com/",
		Description: "INVALID",
		Timestamp:   time.Now(),
	}
	e.Validate()
	if e.Valid == true {
		t.Errorf("error validating event tracker. expected: %t, got: %t\n",
			false, e.Valid)
	}
}
