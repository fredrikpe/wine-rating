package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wine_rating/internal/db"
	"wine_rating/internal/vivino"
)

type MatchRequest struct {
	Name     string `json:"name"`
	Producer string `json:"producer"`
	Year     *int   `json:"year,omitempty"`
}

func MatchHandler(db *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Producer == "" {
			http.Error(w, "Missing 'name' or 'producer'", http.StatusBadRequest)
			return
		}

		match, err := vivino.FindMatch(db, req.Name, req.Producer, req.Year)
		if err != nil {
			http.Error(w, fmt.Sprintf("FindMatch error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

func EnrichHandler(db *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		excel, err := readUploadedExcel(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := enrichExcelWithVivino(db, excel); err != nil {
			http.Error(w, fmt.Sprintf("processing failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=\"enriched.xlsx\"")
		excel.Write(w)
	}
}
