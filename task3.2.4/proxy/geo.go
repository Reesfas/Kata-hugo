package main

import (
	"context"
	"encoding/json"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"log"
	"net/http"
)

type Address struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
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

type Geo struct {
	Responder
}

type Geocoder interface {
	Search(w http.ResponseWriter, r *http.Request)
	geocodeAddress(w http.ResponseWriter, r *http.Request)
}

func NewGeo() Geocoder {
	return &Geo{
		NewResponder(nil, nil),
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
	var request SearchRequest
	var response SearchResponse
	g.Responder.OutputJSON(w, request)
	/*err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}*/
	cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "52a132510af242610a33fea8352874a271dbfebc",
		SecretKeyValue: "c79449d80f97afab795c7f1eda5a746d32391831"}))
	addresses, err := cleanApi.Address(context.Background(), request.Query)
	if err != nil {
		g.ErrorBadRequest(w, err)
		//http.Error(w, "Dadata API is not available", http.StatusInternalServerError)
		return
	}
	log.Println(addresses[0])
	response.Addresses = []*Address{{addresses[0].GeoLat, addresses[0].GeoLon}}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		g.ErrorInternal(w, err)
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// geocodeAddress функция обработки запроса геокодирования.
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
func (g *Geo) geocodeAddress(w http.ResponseWriter, r *http.Request) {
	var requestData GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := &GeocodeResponse{
		Addresses: []*Address{
			{
				Lat: requestData.Lat,
				Lon: requestData.Lon,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
