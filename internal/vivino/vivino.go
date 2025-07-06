package vivino

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"wine_rating/internal/db"
	"wine_rating/internal/similarity"

	"golang.org/x/text/unicode/norm"
)

const RATINGS_COUNT_THRESHOLD = 75

type Match struct {
	Url            string   `json:"url"`
	ExactVintage   bool     `json:"exact_vintage"`
	RatingsAverage *float64 `json:"ratings_average,omitempty"`
	Similarity     float64  `json:"similarity"`
}

func FindMatch(store *db.Store, query string) (Match, error) {
	query, year := parseQuery(query)
	wines, err := getVivinoHits(store, query)
	if err != nil {
		return Match{}, fmt.Errorf("getVivinoHits failed: %w", err)
	}

	best, confidence := bestMatch(wines, query)

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
	wines, updatedAt, err := db.GetVivinoQuery(query)
	if err != nil {
		return nil, fmt.Errorf("get query failed: %w", err)
	}
	if len(wines) > 0 && !updatedAt.Before(time.Now().AddDate(0, 0, -30)) {
		log.Printf("cached query found: %s", query)
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

func Irrelevant(similarity float64) bool {
	return similarity < 0.5
}

func QuiteCertain(similarity float64) bool {
	return similarity > 0.75
}

type NameProducer struct {
	Name     string
	Producer string
}

func WineSimilarity(a, b NameProducer) float64 {
	return similarity.Similarity(
		normalizeQuery(a.Name+" "+a.Producer),
		normalizeQuery(b.Name+" "+b.Producer),
	)
}

func bestMatch(hits []db.VivinoWineDbo, query string) (db.VivinoWineDbo, float64) {
	var best db.VivinoWineDbo
	sim := 0.0
	for _, hit := range hits {
		s := WineSimilarity(
			NameProducer{Name: hit.Name, Producer: hit.Producer},
			NameProducer{Name: query},
		)
		if s > sim {
			best = hit
			sim = s
		}
	}
	return best, sim
}

func parseQuery(q string) (string, *int) {
	year := closestPastYear(extractValidYears(q))

	return normalizeQuery(q), year
}

func extractValidYears(s string) []int {
	re := regexp.MustCompile(`\b(\d{4})\b`)
	matches := re.FindAllString(s, -1)
	currentYear := time.Now().Year()

	var years []int
	for _, m := range matches {
		y, err := strconv.Atoi(m)
		if err != nil || y < 1700 || y > currentYear {
			continue
		}
		years = append(years, y)
	}
	return years
}

func closestPastYear(years []int) *int {
	if len(years) == 0 {
		return nil
	}
	currentYear := time.Now().Year()
	var closest *int
	minDiff := currentYear + 1

	for _, y := range years {
		diff := currentYear - y
		if diff < minDiff {
			minDiff = diff
			val := y
			closest = &val
		}
	}
	return closest
}

func normalizeQuery(q string) string {
	return strings.Join(SortedUnique(stripNumberWords(Normalize(q))), " ")
}

func removeDiacritics(input string) string {
	// Normalize to NFD (decomposed form: é → e +  ́ )
	t := norm.NFD.String(input)

	// Filter out all non-spacing marks (diacritics)
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue // skip diacritic
		}
		b.WriteRune(r)
	}
	return b.String()
}

func stripNumberWords(s string) string {
	numberRe := regexp.MustCompile(`^\d+$`)
	priceRe := regexp.MustCompile(`^\d+,-$`)

	var result []string
	for f := range strings.FieldsSeq(s) {
		if !numberRe.MatchString(f) && !priceRe.MatchString(f) {
			result = append(result, f)
		}
	}
	return strings.Join(result, " ")
}

func Normalize(s string) string {
	return removeDiacritics(strings.ToLower(s))
}

func SortedUnique(s string) []string {
	seen := make(map[string]bool)
	var tokens []string

	for word := range strings.FieldsSeq(s) {
		if !seen[word] {
			seen[word] = true
			tokens = append(tokens, word)
		}
	}

	sort.Strings(tokens)
	return tokens
}
