package repository

import (
	"database/sql"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type GeocoderRepository interface {
	Search(query string) ([]string, error)
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
