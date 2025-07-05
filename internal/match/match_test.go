package match

import (
	"testing"
)

func TestConfidence(t *testing.T) {
	c := Confidence(
		Wine{
			Name:     "Galway Vintage Shiraz",
			Producer: "Yalumba",
			Region:   "South Australia",
			Country:  "au",
		},
		Wine{
			Name:     "Yalumba Galway Vintage Shiraz 2022",
			Producer: "Yalumba Winery",
			Region:   "",
			Country:  "",
		},
	)

	if c < 0.7 {
		t.Fatalf("confidence too low")
	}
}
