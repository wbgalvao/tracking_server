package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wbgalvao/tracking_server/cassandra"

	"github.com/gocql/gocql"
	"github.com/wbgalvao/tracking_server/model"
)

func TestHealthcheck(t *testing.T) {
	var app App
	req, err := http.NewRequest("GET", "localhost:8080/", nil)
	if err != nil {
		t.Errorf("error creating test request: %v\n", err)
	}
	rec := httptest.NewRecorder()
	app.Healthcheck(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error reading Healthcheck test request reponse: %v\n", err)
	}
	var r Response
	json.Unmarshal(b, &r)
	if r.Code != http.StatusOK || r.Message != "OK" {
		t.Errorf("invalid server answer. expected %d and %q, got: %d and %q\n",
			http.StatusOK, "OK", r.Code, r.Message)
	}
}

func TestSendEvent(t *testing.T) {
	cluster := cassandra.Cluster{Hosts: "localhost", Keyspace: "tracker", ProtocolVersion: 4}
	session, err := cassandra.CreateSession(cluster.Create())
	if err != nil {
		t.Errorf("unable to establish connection session with test database: %v\n", err)
	}
	defer session.Close()
	app := App{Session: session, Table: "event"}
	_ = cassandra.CreateEventTable(session, cluster.Keyspace, app.Table)
	e := model.Event{
		ID:          gocql.TimeUUID(),
		Username:    "John Doe",
		Target:      "https://www.google.com/",
		Description: "VIEW",
		Valid:       true,
		Timestamp:   time.Now(),
	}
	eJSON, err := json.Marshal(e)
	if err != nil {
		t.Errorf("unable to marshal test event: %v\n", err)
	}
	req, err := http.NewRequest("POST", "localhost:8080/track", bytes.NewBuffer(eJSON))
	if err != nil {
		t.Errorf("unabe to create test request: %v\n", err)
	}
	rec := httptest.NewRecorder()
	app.TrackEvent(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error reading SendEvent test request response: %v\n", err)
	}
	var r Response
	json.Unmarshal(b, &r)
	if r.Code != http.StatusOK || r.Message != "OK" {
		t.Errorf("invalid server answer. expected %d and %q, got %q and %q\n",
			http.StatusOK, "OK", r.Code, r.Message)
	}
	_ = cassandra.DropTable(session, app.Table)
}
