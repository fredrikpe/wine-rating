package match

import (
	"log"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"wine_rating/internal/levenshtein"

	"golang.org/x/text/unicode/norm"
)

var yearRegex = regexp.MustCompile(`\b\d{4}\b`)

type Wine struct {
	Name     string
	Producer string
	Country  string
	Region   string
}

type Similarity struct {
	Name     float64
	Producer float64
}

func WineSimilarity(a, b Wine) Similarity {
	return Similarity{
		Name:     mongeElkanSimilarity(stripYearAndProducerWords(a), stripYearAndProducerWords(b)),
		Producer: mongeElkanSimilarity(a.Producer, b.Producer),
	}
}

func HighEnough(c float64) bool {
	return c > 0.75
}

func Confidence(a, b Wine) float64 {
	d := WineSimilarity(a, b)

	const (
		wProducer = 0.5
		wName     = 0.5
	)

	c := d.Producer*wProducer + d.Name*wName

	if true {
		log.Printf(`DEBUG: Confidence match
	  Wine A: %-30s
	  Wine B: %-30s

	  Distances:
	    Name:     %.2f (%s ↔ %s)
	    Producer: %.2f (%s ↔ %s)

	  Final Confidence: %.2f
	`,
			a.Name, b.Name,
			d.Name, stripYearAndProducerWords(a), stripYearAndProducerWords(b),
			d.Producer, a.Producer, b.Producer,
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

func mongeElkanSimilarity(a, b string) float64 {
	simAB := mongeElkanSimilarityOneWay(a, b)
	simBA := mongeElkanSimilarityOneWay(b, a)
	return (simAB + simBA) / 2
}

func mongeElkanSimilarityOneWay(a, b string) float64 {
	toksA := strings.Fields(removeDiacritics(strings.ToLower(a)))
	toksB := strings.Fields(removeDiacritics(strings.ToLower(b)))

	var total float64
	for _, tokA := range toksA {
		maxSim := 0.0
		for _, tokB := range toksB {
			sim := 1.0 - levenshtein.NormalizedDistance(tokA, tokB)
			if sim > maxSim {
				maxSim = sim
			}
		}
		total += maxSim
	}
	if len(toksA) == 0 {
		return 0.0
	}
	return total / float64(len(toksA))
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

func stripYearAndProducerWords(w Wine) string {
	name := strings.ToLower(w.Name)
	producer := strings.ToLower(w.Producer)

	// Remove 4-digit year
	name = yearRegex.ReplaceAllString(name, "")

	// Remove producer words
	prodWords := strings.FieldsSeq(producer)
	for word := range prodWords {
		name = strings.ReplaceAll(name, word, "")
	}

	// Normalize spacing
	return strings.Join(strings.Fields(name), " ")
}
