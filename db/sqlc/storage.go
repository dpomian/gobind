package db

import "database/sql"

type Storage interface {
	Querier
}

type SQLStorage struct {
	database *sql.DB
	*Queries
}

func NewStorage(database *sql.DB) Storage {
	return &SQLStorage{
		database: database,
		Queries:  New(database),
	}
}
