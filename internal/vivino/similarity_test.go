package vivino

import (
	"testing"
)

func TestSimilarityGalwayOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Galway Vintage Shiraz",
			Producer: "Yalumba",
		},
		NameProducer{
			Name:     "Yalumba Galway Vintage Shiraz 2022",
			Producer: "Yalumba Winery",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityPironOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Montagne Saint-Émilion",
			Producer: "Château Piron",
		},
		NameProducer{
			Name:     "Ch. Piron Montagne Saint-Émilion 2019",
			Producer: "Fressineau Ch. Piron",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityFaiveley(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Bourgogne Chardonnay",
			Producer: "Domaine Faiveley",
		},
		NameProducer{
			Name:     "Faiveley Bourgogne Chardonnay 2022",
			Producer: "Dom. Faiveley",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityGardiesOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Mas Las Cabes Blanc",
			Producer: "Gardiés",
		},
		NameProducer{
			Name:     "Dom. Gardies Mas Las Cabes Garance 2023",
			Producer: "Dom. Gardiés",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilaritySemeliOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Mantinia Moschofilero",
			Producer: "Seméli",
		},
		NameProducer{
			Name:     "Seméli Mantinia Moschofilero 2023",
			Producer: "SEMELI WINERY S.A.",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityZafeirakisOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Limniona Young Vineyards",
			Producer: "Κτήμα Ζαφειράκη (Domaine Zafeirakis)",
		},
		NameProducer{
			Name:     "Zafeirakis Limniona Young Vineyards 2023",
			Producer: "Zafeirakis",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityVibertiOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Nebbiolo Langhe",
			Producer: "Viberti Giovanni",
		},
		NameProducer{
			Name:     "Viberti Giovanni Langhe Nebbiolo 2022",
			Producer: "Viberti Giovanni",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityDemiereOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Lysandre Cuvée",
			Producer: "Champagne Demière",
		},
		NameProducer{
			Name:     "Demière Lysandre 2009",
			Producer: "Champagne Demière",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilaritySimpsonsOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Chalklands Classic Cuvée",
			Producer: "Simpsons",
		},
		NameProducer{
			Name:     "Simpsons Chalklands Classic Cuvée NV",
			Producer: "Simpsons Wine Estate",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityOlimBaudaOk(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Centive Moscato d'Asti",
			Producer: "Olim Bauda",
		},
		NameProducer{
			Name:     "Olim Bauda Moscato d'Asti Centive 2024",
			Producer: "Ten. Olim Bauda",
		},
	)

	if !QuiteCertain(c) {
		t.Fatalf("confidence too low")
	}
}

func TestSimilarityGardWrong(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Bourgogne Pinot Noir",
			Producer: "Olivier Gard",
		},
		NameProducer{
			Name:     "Lien Gård Eplemost Rubinstep ",
			Producer: "Lien Gård",
		},
	)

	if QuiteCertain(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityRexHillWrong(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Sunny Mountain Vineyard Pinot noir",
			Producer: "Rex Hill",
		},
		NameProducer{
			Name:     "Mikkeller Sunny Shandy",
			Producer: "Mikkeller",
		},
	)

	if QuiteCertain(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityValdesayWrong(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Bourgogne Chardonnay",
			Producer: "Maison Valdesay",
		},
		NameProducer{
			Name:     "Maison Vignoud Bourgogne Chardonnay 2023",
			Producer: "Oslo Wine Agency As",
		},
	)

	if QuiteCertain(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityBrialWrong(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Les Camines Blanc",
			Producer: "Dom Brial",
		},
		NameProducer{
			Name:     "Dom Brial Las Coumeilles 2023",
			Producer: "Dom. Brial",
		},
	)

	if QuiteCertain(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityTroupisWrong(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Moschofilero",
			Producer: "Troupis Winery",
		},
		NameProducer{
			Name:     "Novus A Priori Mantineia Moschofilero 2024",
			Producer: "Novus Winery",
		},
	)

	if QuiteCertain(c) {
		t.Fatalf("confidence too high")
	}
}

func TestSimilarityBarcelonaWrong(t *testing.T) {
	c := WineSimilarity(
		NameProducer{
			Name:     "Cava Brut",
			Producer: "FC Barcelona",
		},
		NameProducer{
			Name: "Hola! Cava Brut Barcelona Brands",
		},
	)

	if QuiteCertain(c) {
		t.Fatalf("confidence too high")
	}
}
