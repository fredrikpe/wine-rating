package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"wine_rating/internal/db"
	"wine_rating/internal/vivino"
)

type MatchRequest struct {
	Name     string `json:"name"`
	Producer string `json:"producer"`
	Year     *int   `json:"year,omitempty"`
}

func (req MatchRequest) toQuery() string {
	var queryParts []string
	if req.Year != nil {
		queryParts = append(queryParts, strconv.Itoa(*req.Year))
	}
	queryParts = append(queryParts, req.Name, req.Producer)
	return strings.Join(queryParts, " ")
}

func WithCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
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

		match, err := vivino.FindMatch(db, req.toQuery())
		if err != nil {
			http.Error(w, fmt.Sprintf("FindMatch error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(match)
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
		if err := excel.Write(w); err != nil {
			log.Printf("failed to write Excel file: %v", err)
			http.Error(w, "Failed to write Excel file", http.StatusInternalServerError)
		}
	}
}
