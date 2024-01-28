package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dpomian/gobind/api"
	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/dpomian/gobind/utils"
	_ "github.com/lib/pq"
)

var storage db.Storage

func init() {

}

func main() {
	startGinServer(":5050")
}

func startGinServer(address string) {
	config, err := utils.LoadConfig("")
	if err != nil {
		log.Fatal("cannot load config")
	}

	fmt.Println("config:", config)

	database, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect do db:", err)
	}

	storage = db.NewStorage(database)

	server, err := api.NewServer(config, storage)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(address)
}
