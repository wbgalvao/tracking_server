package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"

	"github.com/wbgalvao/tracking_server/cassandra"
	"github.com/wbgalvao/tracking_server/handler"
)

func main() {

	hosts := os.Getenv("CASSANDRA_HOST")
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	table := os.Getenv("CASSANDRA_EVENT_TABLE")

	var app handler.App

	cluster := cassandra.Cluster{Hosts: hosts, Keyspace: keyspace}
	session, err := cassandra.CreateSession(cluster.Create())
	if err != nil {
		log.Fatalf("unable to create connection session with Cassandra cluster: %v\n", err)
	}

	err = cassandra.CreateEventTable(session, cluster.Keyspace, table)
	if err != nil {
		log.Fatalf("cannot create table %q in keyspace %q: %v\n", table, cluster.Keyspace, err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", app.Healthcheck).Methods("GET")
	router.HandleFunc("/track", app.TrackEvent).Methods("POST")

	app.Initialize(router, session, table)
	app.Run(":8080")

}
