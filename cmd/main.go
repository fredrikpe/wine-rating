package main

import (
	"database/sql"
	"log"
	"net/http"

	"wine_rating/internal"
	"wine_rating/internal/db"
)

func main() {
	dbSqlite, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer dbSqlite.Close()

	err = db.RunMigrations(dbSqlite)
	if err != nil {
		log.Fatal(err)
	}
	db := db.NewDb(dbSqlite)

	http.HandleFunc("/enrich", internal.EnrichHandler(db))

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
