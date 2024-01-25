package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
)

var dbSource = "postgres://postgres:" + os.Getenv("DKR_POSTGRES_PWD") + "@localhost:5444/binder?sslmode=disable"
var testQueries *Queries

func TestMain(m *testing.M) {
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
