package main

import (
	"log"
	"net/http"

	"wine_rating/internal"
)

func main() {
	http.HandleFunc("/enrich", internal.EnrichHandler)

	log.Println("Server started on :8080")

	http.ListenAndServe(":8080", nil)
}
