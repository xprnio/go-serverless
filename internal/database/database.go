package database

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

type Database struct {
	connection *sql.DB
}

func New(name string) (*Database, error) {
	connection, err := sql.Open("sqlite", name)
	if err != nil {
		return nil, err
	}

	db := &Database{connection}
	if err := db.migrate(); err != nil {
		connection.Close()
		return nil, err
	}

	return db, nil
}

func (d *Database) migrate() error {
	queries := []string{
		"create table if not exists functions (" +
			"id text primary key," +
			"image text," +
			"environment text not null" +
			")",
		"create table if not exists routes (" +
			"id text primary key," +
			"path text unique," +
			"function_id text references functions(id)" +
			")",
	}

	for _, query := range queries {
		if _, err := d.connection.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
