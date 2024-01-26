package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/dpomian/gobind/api"
	db "github.com/dpomian/gobind/db/sqlc"
	_ "github.com/lib/pq"
)

var storage db.Storage

func init() {

}

func main() {
	startGinServer(":5050")
}

func startGinServer(address string) {
	var dbSource = os.Getenv("BINDER_DB_SOURCE")
	var dbDriver = os.Getenv("BINDER_DB_DRIVER")
	database, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect do db:", err)
	}

	storage = db.NewStorage(database)

	server, err := api.NewServer(storage)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(address)
}
