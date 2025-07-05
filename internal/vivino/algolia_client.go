package vivino

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type algoliaResponse struct {
	Hits []VivinoHit `json:"hits"`
}

type VivinoHit struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Winery struct {
		Name   string `json:"name"`
		Region struct {
			Name    string `json:"name"`
			Country string `json:"country"`
		} `json:"region"`
	} `json:"winery"`
	Vintages []struct {
		Id         int          `json:"id"`
		Year       string       `json:"year"`
		Statistics VintageStats `json:"statistics"`
	} `json:"vintages"`
	Statistics WineStats `json:"statistics"`
}

type WineStats struct {
	RatingsAverage float64 `json:"ratings_average"`
	RatingsCount   int     `json:"ratings_count"`
	LabelsCount    int     `json:"labels_count"`
}

type VintageStats struct {
	RatingsAverage float64 `json:"ratings_average"`
	RatingsCount   int     `json:"ratings_count"`
	ReviewsCount   int     `json:"reviews_count"`
	LabelsCount    int     `json:"labels_count"`
}

func algoliaSearch(query string) ([]VivinoHit, error) {
	body, _ := json.Marshal(
		map[string]any{"query": query, "hitsPerPage": 25},
	)

	req, _ := http.NewRequest(
		"POST", "https://9takgwjuxl-dsn.algolia.net/1/indexes/WINES_prod/query",
		bytes.NewBuffer(body),
	)
	req.Header.Set("x-algolia-application-id", "9TAKGWJUXL")
	req.Header.Set("x-algolia-api-key", "60c11b2f1068885161d95ca068d3a6ae")
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Algolia Vivino request - query: %s", body)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vivino request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vivino returned %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var result algoliaResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Hits, nil
}

func decodeVivinoResponse(r io.Reader) ([]VivinoHit, error) {
	var result algoliaResponse
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding vivino response: %w", err)
	}
	return result.Hits, nil
}
