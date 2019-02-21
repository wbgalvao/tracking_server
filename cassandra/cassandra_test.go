package cassandra

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/wbgalvao/tracking_server/model"

	"github.com/gocql/gocql"
)

var (
	cluster         Cluster
	hosts           string
	keyspace        string
	protocolVersion int
	timeout         int
	session         *gocql.Session
	table           string
)

func init() {
	flag.StringVar(&hosts, "hosts", "localhost", "List of hosts addresses for Cassandra cluster (comma separated)")
	flag.StringVar(&keyspace, "keyspace", "tracker", "Cassandra keyspace to use during integration tests")
	flag.IntVar(&timeout, "timeout", 10, "Cluster connection timeout duration (in seconds)")
	flag.IntVar(&protocolVersion, "protoversion", 4, "The version of the native protocol to use")
	flag.StringVar(&table, "table", "event", "Cassandra table to use during integration tests")
	flag.Parse()
}

func setup() {
	cluster = Cluster{
		Hosts:           hosts,
		Keyspace:        keyspace,
		ProtocolVersion: protocolVersion,
		Timeout:         timeout,
	}
	session, _ = CreateSession(cluster.Create())
}

func teardown() {
	_ = DropTable(session, table)
	session.Close()
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func TestCreateCluster(t *testing.T) {
	config := cluster.Create()
	if !reflect.DeepEqual(config.Hosts, strings.Split(hosts, ",")) {
		t.Errorf("cluster created with wrong hosts. expected %q, got %q\n",
			hosts, cluster.Hosts)
	} else if config.Keyspace != keyspace {
		t.Errorf("cluster created with wrong keyspace. expected %q, got %q\n",
			keyspace, cluster.Keyspace)
	}
	if &config == nil {
		t.Errorf("cluster creation returned nil gocql.ClusterConfig\n")
	}
}

func TestCreateSession(t *testing.T) {
	config := cluster.Create()
	_, err := CreateSession(config)
	if err != nil {
		t.Errorf("could not create gocql.Session: %v\n", err)
	}
}

func TestDropTable(t *testing.T) {
	droptestTbl := "droptest"
	q := fmt.Sprintf(`CREATE TABLE %s (pk text, data text, PRIMARY KEY(pk))`, droptestTbl)
	err := session.Query(q).Exec()
	if err != nil {
		t.Errorf("unable to create %q table: %v\n", droptestTbl, err)
	}
	err = DropTable(session, droptestTbl)
	if err != nil {
		t.Errorf("cannot drop table %q: %v\n", droptestTbl, err)
	}
}

func TestCreateEventTable(t *testing.T) {
	session, _ = CreateSession(cluster.Create())
	err := CreateEventTable(session, keyspace, table)
	if err != nil {
		t.Errorf("unable to create table %q: %v\n", table, err)
	}
}

func TestClearDataset(t *testing.T) {
	err := ClearDataset(session, table)
	if err != nil {
		t.Errorf("unable to delete all items from table %s: %v\n", table, err)
	}
}

func TestInsertEvent(t *testing.T) {
	testEvent := model.Event{
		ID:          gocql.TimeUUID(),
		Username:    "John Doe",
		Target:      "https://www.google.com/",
		Description: "CLICK",
		Valid:       true,
		Timestamp:   time.Now(),
	}
	err := InsertEvent(session, table, testEvent)
	if err != nil {
		t.Errorf("unable to insert event in database: %v\n", err)
	}
	_ = ClearDataset(session, table)
}
