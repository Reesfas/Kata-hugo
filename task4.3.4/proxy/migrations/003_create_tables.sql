-- +goose Up
CREATE TABLE IF NOT EXISTS search_history (
id SERIAL PRIMARY KEY,
query TEXT NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS address (
id SERIAL PRIMARY KEY,
address_text TEXT NOT NULL,
geo_lat TEXT NOT NULL,
geo_lon TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS history_search_address (
id SERIAL PRIMARY KEY,
search_history_id INT REFERENCES search_history(id),
address_id INT REFERENCES address(id)
);

