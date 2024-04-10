package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/FeiNiaoBF/simplebank_backend_learn/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// connect to the database
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connection to db:", err)
	}
	if testDB.Ping(); err != nil {
		log.Fatal("invalid connection to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
