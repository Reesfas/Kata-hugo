package repository

import (
	"database/sql"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type GeocoderRepository interface {
	Search(query string) ([]string, error)
	CheckAddressExistence(query string) (bool, error)
	SaveSearchHistory(query string) error
	SaveAddress(addressText string, lat, lon string) (int, error)
}

type GeocodeRep struct {
	db *sql.DB
}

func NewGeoRep(db *sql.DB) *GeocodeRep {
	return &GeocodeRep{db}
}

func (r *GeocodeRep) Search(query string) ([]string, error) {
	var similarAddresses []string

	rows, err := r.db.Query("SELECT address_text FROM address")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var addressText string
		if err = rows.Scan(&addressText); err != nil {
			return nil, err
		}

		distance := levenshtein.RatioForStrings([]rune(query), []rune(addressText), levenshtein.DefaultOptions)

		if distance > 0.7 {
			similarAddresses = append(similarAddresses, addressText)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return similarAddresses, nil
}

// CheckAddressExistence Метод для проверки наличия адреса в базе данных
func (r *GeocodeRep) CheckAddressExistence(query string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM address WHERE address_text ILIKE $1)", "%"+query+"%").Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SaveSearchHistory Метод для сохранения истории поиска
func (r *GeocodeRep) SaveSearchHistory(query string) error {
	_, err := r.db.Exec("INSERT INTO search_history (query) VALUES ($1)", query)
	return err
}

// SaveAddress Метод для сохранения адреса
func (r *GeocodeRep) SaveAddress(addressText string, lat, lon float64) (int, error) {
	var addressID int
	err := r.db.QueryRow("INSERT INTO address (address_text, geo_lat, geo_lon) VALUES ($1, $2, $3) RETURNING id", addressText, lat, lon).Scan(&addressID)
	if err != nil {
		return 0, err
	}
	return addressID, nil
}

// SaveHistorySearchAddress Метод для сохранения связи между адресом и историей поиска
func (r *GeocodeRep) SaveHistorySearchAddress(historyID, addressID int) error {
	_, err := r.db.Exec("INSERT INTO history_search_address (search_history_id, address_id) VALUES ($1, $2)", historyID, addressID)
	return err
}
