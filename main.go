package main

import (
	"database/sql"
	"log"

	"github.com/orkhanrustamli/simplebank/api"
	db "github.com/orkhanrustamli/simplebank/db/sqlc"
	"github.com/orkhanrustamli/simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config file: %v", err)
	}

	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	store := db.NewStore(dbConn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start HTTP server: %v", err)
	}
}
