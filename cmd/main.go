package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
	"wine_rating/internal/db"
	"wine_rating/internal/web"
)

func main() {
	dbSqlite, err := sql.Open("sqlite", "./db.sqlite")
	dbSqlite.SetMaxOpenConns(5)
	dbSqlite.SetMaxIdleConns(5)
	log.Println("Stats:", dbSqlite.Stats())

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := dbSqlite.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	err = db.RunMigrations(dbSqlite)
	if err != nil {
		log.Fatal(err)
	}
	db := db.NewDb(dbSqlite)

	http.Handle("/enrich", web.WithCORS(http.HandlerFunc(web.EnrichHandler(db))))
	http.Handle("/match", web.WithCORS(http.HandlerFunc(web.MatchHandler(db))))
	http.Handle("/match/batch", web.WithCORS(http.HandlerFunc(web.MatchBatchHandler(db))))

	log.Println("Server started on :7661")
	err = http.ListenAndServe(":7661", nil)
	if err != nil {
		log.Fatal(err)
	}
}
