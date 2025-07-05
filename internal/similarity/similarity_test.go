package similarity

import (
	"testing"
)

func TestSimilarityGalwayOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Galway Vintage Shiraz",
			Producer: "Yalumba",
		},
		Wine{
			Name:     "Yalumba Galway Vintage Shiraz 2022",
			Producer: "Yalumba Winery",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityPironOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Montagne Saint-Émilion",
			Producer: "Château Piron",
		},
		Wine{
			Name:     "Ch. Piron Montagne Saint-Émilion 2019",
			Producer: "Fressineau Ch. Piron",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityFaiveley(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Bourgogne Chardonnay",
			Producer: "Domaine Faiveley",
		},
		Wine{
			Name:     "Faiveley Bourgogne Chardonnay 2022",
			Producer: "Dom. Faiveley",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityGardiesOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Mas Las Cabes Blanc",
			Producer: "Gardiés",
		},
		Wine{
			Name:     "Dom. Gardies Mas Las Cabes Garance 2023",
			Producer: "Dom. Gardiés",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilaritySemeliOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Mantinia Moschofilero",
			Producer: "Seméli",
		},
		Wine{
			Name:     "Seméli Mantinia Moschofilero 2023",
			Producer: "SEMELI WINERY S.A.",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityZafeirakisOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Limniona Young Vineyards",
			Producer: "Κτήμα Ζαφειράκη (Domaine Zafeirakis)",
		},
		Wine{
			Name:     "Zafeirakis Limniona Young Vineyards 2023",
			Producer: "Zafeirakis",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Nebbiolo Langhe",
			Producer: "Viberti Giovanni",
		},
		Wine{
			Name:     "Viberti Giovanni Langhe Nebbiolo 2022",
			Producer: "Viberti Giovanni",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityDemiereOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Lysandre Cuvée",
			Producer: "Champagne Demière",
		},
		Wine{
			Name:     "Demière Lysandre 2009",
			Producer: "Champagne Demière",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilaritySimpsonsOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Chalklands Classic Cuvée",
			Producer: "Simpsons",
		},
		Wine{
			Name:     "Simpsons Chalklands Classic Cuvée NV",
			Producer: "Simpsons Wine Estate",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityOlimBaudaOk(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Centive Moscato d'Asti",
			Producer: "Olim Bauda",
		},
		Wine{
			Name:     "Olim Bauda Moscato d'Asti Centive 2024",
			Producer: "Ten. Olim Bauda",
		},
	)

	if !HighEnough(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityGardWrong(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Bourgogne Pinot Noir",
			Producer: "Olivier Gard",
		},
		Wine{
			Name:     "Lien Gård Eplemost Rubinstep ",
			Producer: "Lien Gård",
		},
	)

	if HighEnough(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityRexHillWrong(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Sunny Mountain Vineyard Pinot noir",
			Producer: "Rex Hill",
		},
		Wine{
			Name:     "Mikkeller Sunny Shandy",
			Producer: "Mikkeller",
		},
	)

	if HighEnough(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityValdesayWrong(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Bourgogne Chardonnay",
			Producer: "Maison Valdesay",
		},
		Wine{
			Name:     "Maison Vignoud Bourgogne Chardonnay 2023",
			Producer: "Oslo Wine Agency As",
		},
	)

	if HighEnough(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityBrialWrong(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Les Camines Blanc",
			Producer: "Dom Brial",
		},
		Wine{
			Name:     "Dom Brial Las Coumeilles 2023",
			Producer: "Dom. Brial",
		},
	)

	if HighEnough(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityTroupisWrong(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Moschofilero",
			Producer: "Troupis Winery",
		},
		Wine{
			Name:     "Novus A Priori Mantineia Moschofilero 2024",
			Producer: "Novus Winery",
		},
	)

	if HighEnough(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityBarcelonaWrong(t *testing.T) {
	c := Similarity(
		Wine{
			Name:     "Cava Brut",
			Producer: "FC Barcelona",
		},
		Wine{
			Name:     "Hola! Cava Brut ",
			Producer: "Barcelona Brands",
		},
	)

	if HighEnough(c) {
		t.Fatalf("confidence too high")
	}
}
