package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Store struct {
	Db *sql.DB
}

func NewDb(db *sql.DB) *Store {
	_, err := db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Fatal("Failed to set journal mode:", err)
	}
	_, err = db.Exec(`PRAGMA busy_timeout = 5000;`)
	if err != nil {
		log.Fatal("Failed to set busy_timeout:", err)
	}

	return &Store{Db: db}
}

type VivinoWineDbo struct {
	Id         int
	Name       string
	Producer   string
	Region     string
	Country    string
	Vintages   []VivinoVintageDbo
	Statistics WineStatsDbo
}

type VivinoVintageDbo struct {
	Id           int
	VivinoWineId int
	Year         string
	Statistics   VintageStatsDbo
}

type WineStatsDbo struct {
	RatingsAverage float64
	RatingsCount   int
	LabelsCount    int
}

type VintageStatsDbo struct {
	RatingsAverage float64
	RatingsCount   int
	ReviewsCount   int
	LabelsCount    int
}

func RunMigrations(db *sql.DB) error {
	migrationsDir := getMigrationPath()

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Printf("couldn't find migrations dir: %v", err)
		return nil
	}

	// Sort files to run in order (assuming names are timestamped like 20250705_name.sql)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		path := filepath.Join(migrationsDir, file.Name())
		sqlBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", path, err)
		}

		log.Printf("Running migration: %s", file.Name())
		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("execute migration %s: %w", path, err)
		}
	}

	return nil
}

func (store *Store) UpsertQuery(query string, hits []VivinoWineDbo) error {
	tx, err := store.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed to rollback: %v", err)
			}
			panic(p)
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed to rollback: %v", err)
			}
		}
	}()

	for _, hit := range hits {
		_, err := tx.Exec(`
			INSERT INTO vivino_wine (id, name, producer, region, country, ratings_count, ratings_average, labels_count, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
			ON CONFLICT(id) DO UPDATE SET
				name = excluded.name,
				producer = excluded.producer,
				region = excluded.region,
				country = excluded.country,
				ratings_count = excluded.ratings_count,
				ratings_average = excluded.ratings_average,
				labels_count = excluded.labels_count,
				updated_at = excluded.updated_at
		`, hit.Id, hit.Name, hit.Producer, hit.Region, hit.Country, hit.Statistics.RatingsCount, hit.Statistics.RatingsAverage, hit.Statistics.LabelsCount)

		if err != nil {
			return fmt.Errorf("upsert vivino_wine: %w", err)
		}

		for _, v := range hit.Vintages {
			err = UpsertVivinoVintage(tx, v)
			if err != nil {
				return err
			}
		}
	}

	var queryID int64
	err = tx.QueryRow(`
		INSERT INTO vivino_query (normalized_query, updated_at)
		VALUES (?, CURRENT_TIMESTAMP)
		ON CONFLICT(normalized_query) DO UPDATE
		SET updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`, query).Scan(&queryID)
	if err != nil {
		return fmt.Errorf("upsert vivino_query: %w", err)
	}

	for _, v := range hits {
		_, err := tx.Exec(`
			INSERT OR IGNORE INTO vivino_query_hit (vivino_query_id, vivino_wine_id)
			VALUES (?, ?)`,
			queryID, v.Id)
		if err != nil {
			return fmt.Errorf("insert vivino_query_hit: %w", err)
		}
	}

	return tx.Commit()
}

func UpsertVivinoVintage(tx *sql.Tx, v VivinoVintageDbo) error {
	_, err := tx.Exec(`
		INSERT INTO vivino_vintage (
			id, vivino_wine_id, year, ratings_average, ratings_count, reviews_count, labels_count
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			vivino_wine_id = excluded.vivino_wine_id,
			year = excluded.year,
			ratings_average = excluded.ratings_average,
			ratings_count = excluded.ratings_count,
			reviews_count = excluded.reviews_count,
			labels_count = excluded.labels_count
	`, v.Id, v.VivinoWineId, v.Year, v.Statistics.RatingsAverage, v.Statistics.RatingsCount, v.Statistics.ReviewsCount, v.Statistics.LabelsCount)
	if err != nil {
		return fmt.Errorf("insert vivino_vintage: %w", err)
	}
	return nil
}

func (store *Store) GetVivinoQuery(query string) ([]VivinoWineDbo, time.Time, error) {
	var updatedAt time.Time
	var queryID int64

	err := store.Db.QueryRow(`
		SELECT id, updated_at
		FROM vivino_query
		WHERE normalized_query = ?
	`, query).Scan(&queryID, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, time.Time{}, nil // Not found
	}
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("fetch query metadata: %w", err)
	}

	rows, err := store.Db.Query(`
		SELECT w.id, w.name, w.producer, w.region, w.country, w.ratings_count, w.ratings_average, w.labels_count
		FROM vivino_query_hit qh
		JOIN vivino_wine w ON w.id = qh.vivino_wine_id
		WHERE qh.vivino_query_id = ?
	`, queryID)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("query vivino_wines: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var wines []VivinoWineDbo

	for rows.Next() {
		var wine VivinoWineDbo
		err = rows.Scan(
			&wine.Id,
			&wine.Name,
			&wine.Producer,
			&wine.Region,
			&wine.Country,
			&wine.Statistics.RatingsCount,
			&wine.Statistics.RatingsAverage,
			&wine.Statistics.LabelsCount,
		)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("scan wine: %w", err)
		}

		vintageRows, err := store.Db.Query(`
			SELECT id, vivino_wine_id, year, ratings_average, ratings_count, reviews_count, labels_count
			FROM vivino_vintage
			WHERE vivino_wine_id = ?
			ORDER BY year DESC
		`, wine.Id)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("query vintages for wine %d: %w", wine.Id, err)
		}

		var vintages []VivinoVintageDbo
		for vintageRows.Next() {
			var v VivinoVintageDbo
			err = vintageRows.Scan(
				&v.Id,
				&v.VivinoWineId,
				&v.Year,
				&v.Statistics.RatingsAverage,
				&v.Statistics.RatingsCount, &v.Statistics.ReviewsCount, &v.Statistics.LabelsCount,
			)
			if err != nil {
				if err := vintageRows.Close(); err != nil {
					log.Printf("failed to close rows: %v", err)
				}
				return nil, time.Time{}, fmt.Errorf("scan vintage: %w", err)
			}
			vintages = append(vintages, v)
		}
		if err := vintageRows.Close(); err != nil {
			log.Printf("failed to close vintageRows: %v", err)
		}

		wine.Vintages = vintages
		wines = append(wines, wine)
	}

	if err := rows.Err(); err != nil {
		return nil, time.Time{}, fmt.Errorf("row iteration error: %w", err)
	}

	return wines, updatedAt, nil
}

func getMigrationPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../../db/migrations")
}
