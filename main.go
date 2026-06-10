// Package main is the entry point for the SimpleBank API server.
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/santiagot714/SimpleBank/api"
	db "github.com/santiagot714/SimpleBank/db/sqlc"
	"github.com/santiagot714/SimpleBank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DatabaseURL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
