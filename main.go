package main

import (
	"database/sql"
	"log"

	"github.com/FeiNiaoBF/simplebank_backend_learn/api"
	db "github.com/FeiNiaoBF/simplebank_backend_learn/db/sqlc"
	"github.com/FeiNiaoBF/simplebank_backend_learn/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// connect to the database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to dv:", err)
	}
	if conn.Ping(); err != nil {
		log.Fatal("invalid connection to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
