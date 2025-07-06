package vivino

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestDecodeVivinoResponseFromFile(t *testing.T) {
	f, err := os.Open("testdata/vivino_response.json")
	if err != nil {
		t.Fatalf("failed to open test JSON: %v", err)
	}
	defer f.Close()

	hits, err := decodeVivinoResponse(f)
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	require.Equal(t, 3.7, hits[0].Statistics.RatingsAverage)
}

func TestFindRightTommasiValpolicella(t *testing.T) {
	f, err := os.Open("testdata/tommasi_valpolicella.json")
	if err != nil {
		t.Fatalf("failed to open test JSON: %v", err)
	}
	defer f.Close()

	hits, err := decodeVivinoResponse(f)
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	hit, _ := bestMatch(hitsToDbos(hits), "Valpolicella Tommasi")
	if hit.Id != 1299576 {
		t.Fatalf("Found wrong wine: %v", hit)
	}
}
