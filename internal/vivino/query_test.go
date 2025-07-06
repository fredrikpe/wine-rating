package vivino

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStripNumberWords(t *testing.T) {
	s := stripNumberWords("1 795 hest")
	require.Equal(t, "hest", s)
}

func TestParseQuerySmagOgBehag(t *testing.T) {
	query, year := parseQuery("NV Bruno Paillard Première Cuvée Extra Brut d06/22 pinot noir blend 1 795")
	require.Nil(t, year)
	require.Equal(t, "blend bruno brut cuvee d06/22 extra noir nv paillard pinot premiere", query)

	query, year = parseQuery("2019 Scar of the Sea Methode Ancestrale Topotero Rosé pinot noir 795")
	require.Equal(t, 2019, *year)
	require.Equal(t, "ancestrale methode noir of pinot rose scar sea the topotero", query)

	query, year = parseQuery("2021 Scar of the Sea Methode Ancestrale Topotero Rosé pinot noir 895")
	require.Equal(t, 2021, *year)
	require.Equal(t, "ancestrale methode noir of pinot rose scar sea the topotero", query)

	query, year = parseQuery("2017 Caraccioli Cellars Brut Cuvée chardo blend 1 395")
	require.Equal(t, 2017, *year)
	require.Equal(t, "blend brut caraccioli cellars chardo cuvee", query)

	query, year = parseQuery("2016 Caraccioli Cellars Brut Rosé chardo blend 1 495")
	require.Equal(t, 2016, *year)
	require.Equal(t, "blend brut caraccioli cellars chardo rose", query)

	query, year = parseQuery("PÉT NAT")
	require.Nil(t, year)
	require.Equal(t, "nat pet", query)

	query, year = parseQuery("Frankrike")
	require.Nil(t, year)
	require.Equal(t, "frankrike", query)

	query, year = parseQuery("L16 Domaine des Cavarodes Vin Mousseux La Bulette savagnin blend 1 095")
	require.Nil(t, year)
	require.Equal(t, "blend bulette cavarodes des domaine l16 la mousseux savagnin vin", query)

	query, year = parseQuery("L20 Domaine des Cavarodes Vin Mousseux La Bulette savagnin blend 1 095")
	require.Nil(t, year)
	require.Equal(t, "blend bulette cavarodes des domaine l20 la mousseux savagnin vin", query)

	query, year = parseQuery("2020 Jean-Pierre Robinot Fêtembulles Magnum chenin blanc 1 995")
	require.Equal(t, 2020, *year)
	require.Equal(t, "blanc chenin fetembulles jean-pierre magnum robinot", query)

	query, year = parseQuery("Italia")
	require.Nil(t, year)
	require.Equal(t, "italia", query)
}

func TestParseQueryParkHotell(t *testing.T) {
	query, year := parseQuery("2018 Dönnhoff Oberhäuser Brücke Riesling Grosses Gewachs Versteigerung 1050,- N/A")
	require.Equal(t, 2018, *year)
	require.Equal(t, "brucke donnhoff gewachs grosses n/a oberhauser riesling versteigerung", query)

	query, year = parseQuery("2016 Dönnhoff Felsenberg Riesling Grosses Gewachs 180,- 900,-")
	require.Equal(t, 2016, *year)
	require.Equal(t, "donnhoff felsenberg gewachs grosses riesling", query)

	query, year = parseQuery("2013 Weingut Hermannsberg Steinterrassen Riesling Trocken 150,- 750,-")
	require.Equal(t, 2013, *year)
	require.Equal(t, "hermannsberg riesling steinterrassen trocken weingut", query)

	query, year = parseQuery("2004 Emrich-Schönleber Monziger Frülingsplätzen Riesling Spätlese 135,- 675,-")
	require.Equal(t, 2004, *year)
	require.Equal(t, "emrich-schonleber frulingsplatzen monziger riesling spatlese", query)

	query, year = parseQuery("2017 Weingut Kruger-Rumpf Abtei Riesling Trocken 130,- 650,-")
	require.Equal(t, 2017, *year)
	require.Equal(t, "abtei kruger-rumpf riesling trocken weingut", query)

	query, year = parseQuery("RHEINHESSEN")
	require.Nil(t, year)
	require.Equal(t, "rheinhessen", query)

	query, year = parseQuery("2015 Kai Schätzel Pettenthal Riesling Grosses Gewachs 400,- 2000,-")
	require.Equal(t, 2015, *year)
	require.Equal(t, "gewachs grosses kai pettenthal riesling schatzel", query)

	query, year = parseQuery("2017 Kai Schätzel Hipping Riesling Grosses Gewahs 160,- 800,-")
	require.Equal(t, 2017, *year)
	require.Equal(t, "gewahs grosses hipping kai riesling schatzel", query)

	query, year = parseQuery("2018 Weingut Wittmann Pinot Gris But Colored 140,- 700,-")
	require.Equal(t, 2018, *year)
	require.Equal(t, "but colored gris pinot weingut wittmann", query)

	query, year = parseQuery("2015 Lisa Bunn Hipping Riesling Trocken 140,- 700,-")
	require.Equal(t, 2015, *year)
	require.Equal(t, "bunn hipping lisa riesling trocken", query)
}
