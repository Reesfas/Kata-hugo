-- +goose Up
CREATE TABLE IF NOT EXISTS search_history (
id SERIAL PRIMARY KEY,
query TEXT NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица адресов
CREATE TABLE IF NOT EXISTS address (
id SERIAL PRIMARY KEY,
address_text TEXT NOT NULL,
geo_lat NUMERIC,
geo_lon NUMERIC
);

-- Таблица адресов в истории поиска
CREATE TABLE IF NOT EXISTS history_search_address (
id SERIAL PRIMARY KEY,
search_history_id INT REFERENCES search_history(id),
address_id INT REFERENCES address(id)
);