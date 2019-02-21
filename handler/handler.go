package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wbgalvao/tracking_server/cassandra"

	"github.com/gocql/gocql"

	"github.com/gorilla/mux"

	"github.com/wbgalvao/tracking_server/model"
)

// App represents the server application to handle HTTP/HTTPS requests.
type App struct {
	Router  *mux.Router
	Session *gocql.Session
	Table   string
}

// Initialize gets the app ready for running, instantiating an HTTP router
// and a session to connect with a Cassandra cluster.
func (a *App) Initialize(router *mux.Router, session *gocql.Session, table string) {
	a.Router = router
	a.Session = session
	a.Table = table
}

// Run starts the application web server, serving http in a given address.
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Healthcheck is used as handler to requests whose objective is to check the
// application health. If the request is successful it means the server is ok,
// therefore the handler will respond confirm this status in its response.
func (a *App) Healthcheck(w http.ResponseWriter, r *http.Request) {
	res := Response{Message: "OK", Code: http.StatusOK}
	respondWithJSON(w, http.StatusOK, res)
}

// TrackEvent will read a HTTP request body and try to unmarshall its body
// into an Event structure. In case this operation is successful, the created
// Event will be stored
func (a *App) TrackEvent(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var e model.Event
	err := decoder.Decode(&e)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect json format")
		return
	}
	e.Validate()
	if !e.Valid {
		respondWithError(w, http.StatusBadRequest, "invalid event")
		return
	}
	err = cassandra.InsertEvent(a.Session, a.Table, e)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res := Response{Message: "OK", Code: http.StatusOK}
	respondWithJSON(w, http.StatusOK, res)
}

// Response defines the response structure for the handler functions.
type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
