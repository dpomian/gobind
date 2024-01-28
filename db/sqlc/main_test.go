package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TestData struct {
	userId1 uuid.UUID
}

var testQueries *Queries
var testData TestData

func TestMain(m *testing.M) {
	var dbSource = os.Getenv("BINDER_DB_SOURCE")
	var dbDriver = os.Getenv("BINDER_DB_DRIVER")

	testData = TestData{
		userId1: uuid.New(),
	}

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
