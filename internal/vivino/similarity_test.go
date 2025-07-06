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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if !HighEnough(c) {
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

	if HighEnough(c) {
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

	if HighEnough(c) {
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

	if HighEnough(c) {
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

	if HighEnough(c) {
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

	if HighEnough(c) {
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
			Name:     "Hola! Cava Brut ",
			Producer: "Barcelona Brands",
		},
	)

	if HighEnough(c) {
		t.Fatalf("confidence too high")
	}
}

//2025/07/06 08:28:49 galway 0.9107142857142857
//2025/07/06 08:28:49 piron 0.798731884057971
//2025/07/06 08:28:49 faivel 0.9215686274509803
//2025/07/06 08:28:49 gardies 0.8159722222222221
//2025/07/06 08:28:49 semeli 0.7555555555555555
//2025/07/06 08:28:49 zafei 0.6787280701754386
//2025/07/06 08:28:49 veberti 0.6658119658119658
////2025/07/06 08:28:49 demiere 0.7401960784313726
//2025/07/06 08:28:49 simpsom 0.7307692307692307
//2025/07/06 08:28:49 olim 0.7681818181818182
//2025/07/06 08:28:49 gard 0.4875801282051282
//2025/07/06 08:28:49 rex 0.49722222222222223
//2025/07/06 08:28:49 valdesay 0.560515873015873
//2025/07/06 08:28:49 brial 0.5997023809523809
//2025/07/06 08:28:49 troupis 0.6622807017543859
//2025/07/06 08:28:49 FC b 0.6243589743589744

//2025/07/06 08:29:52 galway 0.875
//2025/07/06 08:29:52 piron 0.7946428571428572
//2025/07/06 08:29:52 faivel 0.8571428571428572
//2025/07/06 08:29:52 gardies 0.8
//2025/07/06 08:29:52 semeli 0.7916666666666667
//2025/07/06 08:29:52 zafei 0.790340909090909
//2025/07/06 08:29:52 veberti 1
//2025/07/06 08:29:52 demiere 0.890625
//2025/07/06 08:29:52 simpsom 0.8041666666666667
//2025/07/06 08:29:52 olim 0.9166666666666666
//2025/07/06 08:29:52 gard 0.49661044973544977
//2025/07/06 08:29:52 rex 0.41138888888888886
//2025/07/06 08:29:52 valdesay 0.5543402777777777
//2025/07/06 08:29:52 brial 0.6771825396825397
//2025/07/06 08:29:52 troupis 0.659970238095238
//2025/07/06 08:29:52 FC b 0.7250000000000001
