package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	var dbSource = os.Getenv("BINDER_DB_SOURCE")
	var dbDriver = os.Getenv("BINDER_DB_DRIVER")

	conn, err := sql.Open(dbDriver, dbSource)

	tx := beginTx(conn)
	defer rollbackTx(tx)

	if err != nil {
		log.Fatal("cannot connect do db:", err)
	}

	testQueries = New(tx)

	os.Exit(m.Run())
}

func beginTx(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("could not begin transaction: %v", err)
	}
	return tx
}

func rollbackTx(tx *sql.Tx) {
	if tx != nil {
		tx.Rollback()
	}
}

func commitTx(tx *sql.Tx) {
	if tx != nil {
		tx.Commit()
	}
}
