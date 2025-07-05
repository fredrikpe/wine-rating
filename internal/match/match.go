package match

import (
	"math"
	"wine_rating/internal/levenshtein"
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
		Name:     levenshtein.NormalizedDistance(a.Name, b.Name),
		Producer: levenshtein.NormalizedDistance(a.Producer, b.Producer),
		Region:   levenshtein.NormalizedDistance(a.Region, b.Region),
		Country:  levenshtein.NormalizedDistance(a.Country, b.Country),
	}
}

func Similarity(a, b Wine) int {
	distance := WineDistance(a, b)

	score := 0
	score += int(math.Round(distance.Name * 5))
	score += int(math.Round(distance.Producer * 5))
	if distance.Region <= 0.4 {
		score += 1
	}
	if distance.Country <= 0.4 {
		score += 1
	}

	return score
}
