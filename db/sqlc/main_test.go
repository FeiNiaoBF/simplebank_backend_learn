package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	// connect to the database
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("cannot connection to db:", err)
	}
	if conn.Ping(); err != nil {
		log.Fatal("invalid connection to db:", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
