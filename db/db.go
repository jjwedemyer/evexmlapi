package data

import (
	"database/sql"
	"log"
)

const (
	DbConnString string = "user=postgres password=postgres dbname=SQO sslmode=disable"
	DbDialect    string = "postgres"
)

type Key struct {
	CharacterId string
	KeyId       string
	Path        string
}

type DB struct {
	*sql.DB
}

func NewDB() *DB {
	db, err := sql.Open(DbDialect, DbConnString)
	if err != nil {
		log.Fatal("Error Opening DB: ", err)
		return nil
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error Pinging DB: ", err)
		return nil
	}
	return &DB{db}
}
