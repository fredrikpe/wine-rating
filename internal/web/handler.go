package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"wine_rating/internal/db"
	"wine_rating/internal/vivino"
)

type MatchRequest struct {
	Query string `json:"query"`
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

		if req.Query == "" {
			http.Error(w, "empty 'query'", http.StatusBadRequest)
			return
		}

		match, err := vivino.FindMatch(db, req.Query)
		if err != nil {
			http.Error(w, fmt.Sprintf("FindMatch error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(match)
	}
}

type MatchBatchResult struct {
	Query string        `json:"query"`
	Match *vivino.Match `json:"match,omitempty"`
	Error string        `json:"error,omitempty"`
}

func MatchBatchHandler(db *db.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requests []MatchRequest
		if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var responses []MatchBatchResult
		for _, req := range requests[:min(len(requests), 50)] {
			if len(req.Query) < 10 {
				responses = append(responses, MatchBatchResult{
					Query: req.Query,
					Error: "invalid query",
				})
				continue
			}
			match, err := vivino.FindMatch(db, req.Query)
			if err != nil {
				responses = append(responses, MatchBatchResult{
					Query: req.Query,
					Error: fmt.Sprintf("%v", err),
				})
				continue
			}
			if !vivino.QuiteCertain(match.Similarity) {
				responses = append(responses, MatchBatchResult{
					Query: req.Query,
					Error: "not found",
				})
				continue
			}
			responses = append(responses, MatchBatchResult{
				Query: req.Query,
				Match: &match,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responses)
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
