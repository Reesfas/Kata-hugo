package repository

import (
	"database/sql"
	"github.com/ekomobile/dadata/v2/api/suggest"
)

type Address struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	AddressText string
}

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

type Connect struct {
	Connect *suggest.Api
}

type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lon string `json:"lng"`
}

type GeocoderRepository interface {
	Search(query string) ([]Address, error)
	CheckAddressExistence(query string) (bool, error)
	SaveSearchHistory(query string) (int, error)
	SaveAddress(addressText string, lat, lon string) (int, error)
	SaveHistorySearchAddress(historyID, addressID int) error
}

type GeocodeRep struct {
	db *sql.DB
}

func NewGeoRep(db *sql.DB) *GeocodeRep {
	return &GeocodeRep{db}
}

func (r *GeocodeRep) Search(query string) ([]Address, error) {
	var similarAddresses []Address
	var address Address
	rows, err := r.db.Query("SELECT * FROM address WHERE levenshtein(address_text, $1) <= LENGTH('$1') * 0.7;", query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&address.AddressText, &address.Lat, &address.Lat); err != nil {
			return nil, err
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
func (r *GeocodeRep) SaveSearchHistory(query string) (int, error) {
	var historyID int
	err := r.db.QueryRow("INSERT INTO search_history (query) VALUES ($1) RETURNING id", query).Scan(&historyID)
	return historyID, err
}

// SaveAddress Метод для сохранения адреса
func (r *GeocodeRep) SaveAddress(addressText string, lat, lon string) (int, error) {
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
