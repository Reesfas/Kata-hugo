package controller

import (
	"encoding/json"
	"hugoproxy-main/task3.3.3/proxy/internal/service"
	"log"
	"net/http"
)

type Geo struct {
	service.Geocoder
	Responder
}

type Geocoderer interface {
	Search(w http.ResponseWriter, r *http.Request)
	GeocodeAddress(w http.ResponseWriter, r *http.Request)
}

func NewGeo(geocoder service.Geocoder, responder Responder) Geocoderer {
	return &Geo{
		geocoder,
		responder,
	}
}

// Search функция обработки запроса поиска адреса.
//
// @summary Поиск адреса по запросу.
// @description Позволяет найти адрес по заданному запросу, используя Dadata API.
// @tags address
// @id searchAddress
// @produce json
// @param request body SearchRequest true "Запрос на поиск адреса"
// @success 200 {object} SearchResponse "Успешный ответ"
// @failure 400 {object} string "Неверный формат запроса"
// @failure 500 {object} string "Dadata API недоступен"
// @router /api/address/search [post]
func (g *Geo) Search(w http.ResponseWriter, r *http.Request) {
	var request service.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	addresses, err := g.Geocoder.SearchService(request)
	if err != nil {
		g.ErrorBadRequest(w, err)
		return
	}
	log.Println(addresses[0])
	response := service.SearchResponse{Addresses: addresses}
	g.Responder.OutputJSON(w, response)
}

// GeocodeAddress функция обработки запроса геокодирования.
//
// @summary Геокодирование координат.
// @description Позволяет выполнить геокодирование по заданным координатам.
// @tags address
// @id geocodeAddress
// @produce json
// @param request body GeocodeRequest true "Запрос на геокодирование"
// @success 200 {object} GeocodeResponse "Успешный ответ"
// @failure 400 {object} string "Неверный формат запроса"
// @router /api/address/geocode [post]
func (g *Geo) GeocodeAddress(w http.ResponseWriter, r *http.Request) {
	var requestData service.GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := g.Geocoder.GeocodeAddressService(requestData.Lat, requestData.Lon)
	if err != nil {
		http.Error(w, "Failed to geocode", http.StatusInternalServerError)
		return
	}

	response := service.GeocodeResponse{Addresses: []*service.Address{address}}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
