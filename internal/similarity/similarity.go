package similarity

import (
	"log"
	"strings"
	"wine_rating/internal/levenshtein"
)

func Similarity(a, b string) float64 {
	sim := mongeElkanSimilarity(a, b)

	if true {
		log.Printf(`DEBUG: Similarity: %.2f (%s ↔ %s)`,
			sim, a, b,
		)
	}

	return sim
}

func mongeElkanSimilarity(a, b string) float64 {
	simAB := mongeElkanSimilarityOneWay(a, b)
	simBA := mongeElkanSimilarityOneWay(b, a)
	return (simAB + simBA) / 2
}

func mongeElkanSimilarityOneWay(a, b string) float64 {
	toksA := strings.Fields(a)
	toksB := strings.Fields(b)

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
