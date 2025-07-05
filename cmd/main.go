package main

import (
	"database/sql"
	"log"
	"net/http"

	"wine_rating/internal"
	"wine_rating/internal/db"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	err := runMigrations()
	if err != nil {
		log.Fatal(err)
	}

	dbSqlite, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer dbSqlite.Close()
	db := db.NewDb(dbSqlite)

	http.HandleFunc("/enrich", internal.EnrichHandler(db))

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func runMigrations() error {
	m, err := migrate.New(
		"file://db/migrations",
		"sqlite3://db.sqlite",
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
