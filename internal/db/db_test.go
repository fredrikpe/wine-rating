package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *Store {
	sqlDB, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(t, err)

	store := NewDb(sqlDB)

	require.NoError(t, RunMigrations(store.Db))

	return store
}

func TestInsertAndGetVivinoQuery(t *testing.T) {
	store := setupTestDB(t)

	query := "gaja sito moresco"
	wine := VivinoWineDbo{
		Id:       1,
		Name:     "Sito Moresco",
		Producer: "Gaja",
		Region:   "Langhe",
		Country:  "Italy",
		Statistics: WineStatsDbo{
			RatingsAverage: 4.1,
			RatingsCount:   30000,
			LabelsCount:    200,
		},
		Vintages: []VivinoVintageDbo{
			{
				Id:           101,
				VivinoWineId: 1,
				Year:         "2020",
				Statistics: VintageStatsDbo{
					RatingsAverage: 4.2,
					RatingsCount:   1500,
					ReviewsCount:   120,
					LabelsCount:    100,
				},
			},
		},
	}

	err := store.UpsertQuery(query, []VivinoWineDbo{wine})
	require.NoError(t, err)

	// Fetch and verify
	got, updatedAt, err := store.GetVivinoQuery(query)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAt)
	require.Len(t, got, 1)
	require.Equal(t, wine.Id, got[0].Id)
	require.Equal(t, wine.Statistics, got[0].Statistics)
	require.Equal(t, wine.Name, got[0].Name)
	require.Equal(t, "2020", got[0].Vintages[0].Year)
	require.Equal(t, 4.2, got[0].Vintages[0].Statistics.RatingsAverage)
}
