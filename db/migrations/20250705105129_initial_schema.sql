CREATE TABLE IF NOT EXISTS vivino_wine (
    id              INTEGER PRIMARY KEY,
    name            TEXT NOT NULL,
    producer        TEXT,
    region          TEXT,
    country         TEXT,
    ratings_average REAL NOT NULL,
    ratings_count   INTEGER NOT NULL,
    labels_count    INTEGER NOT NULL,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vivino_vintage (
    id               INTEGER PRIMARY KEY,
    vivino_wine_id   INTEGER NOT NULL,
    year             TEXT,
    ratings_average  REAL,
    ratings_count    INTEGER,
    reviews_count    INTEGER,
    labels_count     INTEGER,
    FOREIGN KEY (vivino_wine_id) REFERENCES vivino_wine(id) ON DELETE CASCADE,
    UNIQUE (vivino_wine_id, year)
);

CREATE TABLE IF NOT EXISTS vivino_query (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    normalized_query  TEXT NOT NULL UNIQUE,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vivino_query_hit (
    vivino_query_id    INTEGER NOT NULL,
    vivino_wine_id     INTEGER NOT NULL,
    PRIMARY KEY (vivino_query_id, vivino_wine_id),
    FOREIGN KEY (vivino_query_id) REFERENCES vivino_query(id) ON DELETE CASCADE,
    FOREIGN KEY (vivino_wine_id) REFERENCES vivino_wine(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS grape (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    vivino_id     INTEGER NOT NULL UNIQUE,
    name          TEXT NOT NULL,
    synonyms      TEXT
);

