package match

import (
	"math"
	"sort"
	"strings"
	"unicode"
	"wine_rating/internal/levenshtein"

	"golang.org/x/text/unicode/norm"
)

type Wine struct {
	Name     string
	Producer string
	Country  string
	Region   string
}

type Distance struct {
	Name     float64
	Producer float64
	Country  float64
	Region   float64
}

func WineDistance(a, b Wine) Distance {
	return Distance{
		Name:     jaccardLev(a.Name, b.Name),
		Producer: jaccardLev(a.Producer, b.Producer),
		Region:   jaccardLev(a.Region, b.Region),
		Country:  jaccardLev(a.Country, b.Country),
	}
}

func Similarity(a, b Wine) int {
	distance := WineDistance(a, b)

	score := 0
	score += int(math.Round(distance.Name * 5))
	score += int(math.Round(distance.Producer * 10))
	if distance.Region <= 0.4 {
		score += 1
	}
	if distance.Country <= 0.4 {
		score += 1
	}

	return score
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

func jaccardLev(a, b string) float64 {
	ca := removeDiacritics(strings.ToLower(a))
	cb := removeDiacritics(strings.ToLower(b))

	jaccard := jaccardSimilarity(SortedUnique(ca), SortedUnique(cb))
	lev := levenshtein.NormalizedDistance(ca, cb)

	return 0.7*jaccard + 0.3*(1-lev)
}

func jaccardSimilarity(a, b []string) float64 {
	setA := make(map[string]struct{})
	setB := make(map[string]struct{})

	for _, token := range a {
		setA[token] = struct{}{}
	}
	for _, token := range b {
		setB[token] = struct{}{}
	}

	intersection := 0
	union := make(map[string]struct{})

	for token := range setA {
		union[token] = struct{}{}
		if _, ok := setB[token]; ok {
			intersection++
		}
	}
	for token := range setB {
		union[token] = struct{}{}
	}

	if len(union) == 0 {
		return 0.0
	}
	return float64(intersection) / float64(len(union))
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
