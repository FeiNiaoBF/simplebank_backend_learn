package main

import (
	"database/sql"
	"log"

	"github.com/FeiNiaoBF/simplebank_backend_learn/api"
	db "github.com/FeiNiaoBF/simplebank_backend_learn/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// connect to the database
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to dv:", err)
	}
	if conn.Ping(); err != nil {
		log.Fatal("invalid connection to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
