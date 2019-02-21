package cassandra

import (
	"errors"
	"fmt"
	"time"

	"github.com/wbgalvao/tracking_server/model"

	"github.com/gocql/gocql"
)

var (
	// ErrInvalidEvent is returned by InsertEvent when there is an attempt to
	// insert an invalid model.Tracker into the database.
	ErrInvalidEvent = errors.New("cassandra: cannot insert invalid event into the database")
)

// Cluster defines the basic configurations for creation and management of
// Apache Cassandra databases clusters.
//	- Hosts defines the addresses of Cassandra's cluster nodes
//	- Keyspace defines the cluster keyspace (name, identification, alias...)
//	- ProtocolVersion defines the version of Cassandra communication protocol
type Cluster struct {
	Hosts           string
	Keyspace        string
	ProtocolVersion int
	Timeout         int
}

// CreateCluster wraps the gocql.NewCluster function. It returns a
// *gocql.ClusterConfig which can be used to create sessions to connect
// to existing Cassandra clusters.
func (c Cluster) Create() gocql.ClusterConfig {
	config := gocql.NewCluster(c.Hosts)
	config.Keyspace = c.Keyspace
	if c.ProtocolVersion == 0 {
		c.ProtocolVersion = 4
	}
	config.ProtoVersion = c.ProtocolVersion
	if c.Timeout == 0 {
		c.Timeout = 10
	}
	config.Timeout = time.Duration(c.Timeout) * time.Second
	return *config
}

// CreateSession wraps the gocql.NewSession function. It returns a
// *gocql.Session, which can be used to run queries in an existing Cassandra
// cluster.
func CreateSession(config gocql.ClusterConfig) (*gocql.Session, error) {
	return gocql.NewSession(config)
}

// DropTable will execute a DROP cql expression in a given table.
func DropTable(session *gocql.Session, table string) error {
	q := fmt.Sprintf(`DROP TABLE %s`, table)
	return session.Query(q).Exec()
}

// CreateTable creates a table in the Cassandra cluster which the given session
// created a connection with.
func CreateEventTable(session *gocql.Session, keyspace, table string) error {
	q := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s(
			id             uuid,
			username       text,
			target         text,
			description    text,
			timestamp      timestamp,
		PRIMARY KEY(id)
		)`, table)
	return session.Query(q).Exec()
}

// ClearDataset will truncate a given table, removing all items presented in it.
func ClearDataset(session *gocql.Session, table string) error {
	q := fmt.Sprintf("TRUNCATE TABLE %s", table)
	return session.Query(q).Exec()
}

// InsertEvent uses a given session to insert the data presented in a given
// event tracker into an also given table. If the event is not valid, returns
// an ErrInvalidEvent error.
func InsertEvent(session *gocql.Session, table string, event model.Event) error {
	if !event.Valid {
		return ErrInvalidEvent
	}
	q := fmt.Sprintf(
		`INSERT into %s (id, username, target, description, timestamp)
		 VALUES (?, ?, ?, ?, ?)`, table)
	return session.Query(q, event.ID, event.Username,
		event.Target, event.Description, event.Timestamp).Exec()
}
