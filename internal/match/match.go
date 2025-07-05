package match

import (
	"log"
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
		Name:     jaccardLevDistance(a.Name, b.Name),
		Producer: jaccardLevDistance(a.Producer, b.Producer),
		Region:   jaccardLevDistance(a.Region, b.Region),
		Country:  jaccardLevDistance(a.Country, b.Country),
	}
}

func Confidence(a, b Wine) float64 {
	d := WineDistance(a, b)

	nameSim := 1 - d.Name
	producerSim := 1 - d.Producer
	regionSim := 1 - d.Region
	countrySim := 1 - d.Country

	const (
		wProducer = 0.50
		wName     = 0.3
		wRegion   = 0.05
		wCountry  = 0.15
	)

	c := producerSim*wProducer +
		nameSim*wName +
		regionSim*wRegion +
		countrySim*wCountry
	if false {
		log.Printf(`DEBUG: Confidence match
	  Wine A: %-30s
	  Wine B: %-30s

	  Distances:
	    Name:     %.2f (%s ↔ %s)
	    Producer: %.2f (%s ↔ %s)
	    Region:   %.2f (%s ↔ %s)
	    Country:  %.2f (%s ↔ %s)

	  Final Confidence: %.2f
	`,
			a.Name, b.Name,
			d.Name, a.Name, b.Name,
			d.Producer, a.Producer, b.Producer,
			d.Region, a.Region, b.Region,
			d.Country, a.Country, b.Country,
			c,
		)
	}

	return c
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

func jaccardLevDistance(a, b string) float64 {
	ca := removeDiacritics(strings.ToLower(a))
	cb := removeDiacritics(strings.ToLower(b))

	jaccard := jaccardSimilarity(SortedUnique(ca), SortedUnique(cb))
	lev := levenshtein.NormalizedDistance(ca, cb)

	return 0.7*(1-jaccard) + 0.3*lev
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
