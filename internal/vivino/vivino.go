package vivino

import (
	"fmt"
	"log"
	"strings"
	"wine_rating/internal/match"
)

func FindVivinoMatch(name, producer string) (VivinoHit, error) {
	hits, err := getVivinoHits(normalizeQuery(name, producer))
	if err != nil {
		return VivinoHit{}, err
	}

	best, distance := bestMatch(hits, match.Wine{
		Name:     name,
		Producer: producer,
	})
	log.Printf("DEBUG: Best match id=%d distance=%+v",
		best.Id, distance)

	return best, nil
}

func Url(id int) string {
	return fmt.Sprintf("https://vivino.com/w/%d", id)
}

func bestMatch(hits []VivinoHit, wine match.Wine) (VivinoHit, match.Distance) {
	max := VivinoHit{}
	maxScore := 0
	for _, hit := range hits {
		score := match.Similarity(
			match.Wine{
				Name:     hit.Name,
				Producer: hit.Winery.Name,
				Region:   hit.Region.Name,
				Country:  hit.Region.Country,
			},
			wine,
		)
		if score > maxScore {
			max = hit
			maxScore = score
		}
	}
	return max, match.WineDistance(
		match.Wine{
			Name:     max.Name,
			Producer: max.Winery.Name,
			Region:   max.Region.Name,
			Country:  max.Region.Country,
		},
		wine,
	)
}

func normalizeQuery(name, producer string) string {
	return strings.Join(match.SortedUnique(match.Normalize(name+" "+producer)), " ")
}
