package vivino

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"wine_rating/internal/db"
	"wine_rating/internal/similarity"
)

const RATINGS_COUNT_THRESHOLD = 75

type Match struct {
	Url            string
	ExactVintage   bool
	RatingsAverage *float64
	Similarity     float64
}

func FindMatch(store *db.Store, name, producer string, year *int) (Match, error) {
	wines, err := getVivinoHits(store, normalizeQuery(name, producer))
	if err != nil {
		return Match{}, fmt.Errorf("getVivinoHits failed: %w", err)
	}

	best, confidence := bestMatch(wines, similarity.Wine{
		Name:     name,
		Producer: producer,
	})

	var ratingsAverage *float64
	exactVintage := false

	if year != nil {
		for _, vintage := range best.Vintages {
			if strconv.Itoa(*year) == vintage.Year && vintage.Statistics.RatingsCount > RATINGS_COUNT_THRESHOLD {
				ratingsAverage = &vintage.Statistics.RatingsAverage
				exactVintage = true
				break
			}
		}
	}
	if ratingsAverage == nil && best.Statistics.RatingsCount > RATINGS_COUNT_THRESHOLD {
		ratingsAverage = &best.Statistics.RatingsAverage
	}

	return Match{
		Url:            Url(best.Id),
		ExactVintage:   exactVintage,
		RatingsAverage: ratingsAverage,
		Similarity:     confidence,
	}, nil
}

func Url(id int) string {
	return fmt.Sprintf("https://vivino.com/w/%d", id)
}

func getVivinoHits(db *db.Store, query string) ([]db.VivinoWineDbo, error) {
	wines, _, err := db.GetVivinoQuery(query)
	if err != nil {
		return nil, fmt.Errorf("get query failed: %w", err)
	}
	if len(wines) > 0 {
		log.Println("Returning wines from db")
		return wines, nil
	}
	hits, err := algoliaSearch(query)
	if err != nil {
		return nil, fmt.Errorf("get query failed: %w", err)
	}
	dbos := hitsToDbos(hits)

	err = db.UpsertQuery(query, dbos)
	if err != nil {
		return nil, fmt.Errorf("get query failed: %w", err)
	}

	return dbos, nil
}

func hitsToDbos(hits []VivinoHit) []db.VivinoWineDbo {
	var result []db.VivinoWineDbo
	for _, hit := range hits {
		result = append(result, hitToDbo(hit))
	}
	return result
}

func hitToDbo(hit VivinoHit) db.VivinoWineDbo {
	wine := db.VivinoWineDbo{
		Id:         hit.Id,
		Name:       hit.Name,
		Producer:   hit.Winery.Name,
		Region:     hit.Winery.Region.Name,
		Country:    hit.Winery.Region.Country,
		Statistics: db.WineStatsDbo(hit.Statistics),
	}

	for _, v := range hit.Vintages {
		wine.Vintages = append(wine.Vintages, db.VivinoVintageDbo{
			Id:           v.Id,
			VivinoWineId: hit.Id,
			Year:         v.Year,
			Statistics:   db.VintageStatsDbo(v.Statistics),
		})
	}

	return wine
}

func bestMatch(hits []db.VivinoWineDbo, wine similarity.Wine) (db.VivinoWineDbo, float64) {
	var best db.VivinoWineDbo
	confidence := 0.0
	for _, hit := range hits {
		c := similarity.Similarity(
			similarity.Wine{
				Name:     hit.Name,
				Producer: hit.Producer,
				Region:   hit.Region,
				Country:  hit.Country,
			},
			wine,
		)
		if c > confidence {
			best = hit
			confidence = c
		}
	}
	return best, confidence
}

func normalizeQuery(name, producer string) string {
	return strings.Join(similarity.SortedUnique(similarity.Normalize(name+" "+producer)), " ")
}
