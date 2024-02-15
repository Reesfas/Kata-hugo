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
func Search(w http.ResponseWriter, r *http.Request) {
	var request SearchRequest
	var response SearchResponse
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "11cb4969967b7e68ab87b57258372aefec0eb6ac",
		SecretKeyValue: "3461265109aaa28b20523e1b4dfb4d36e475fc9f"}))
	addresses, err1 := cleanApi.Address(context.Background(), request.Query)
	if err1 != nil {
		http.Error(w, "Dadata API is not available", http.StatusInternalServerError)
		return
	}
	log.Println(addresses[0])
	response.Addresses = []*Address{{addresses[0].GeoLat, addresses[0].GeoLon}}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
func geocodeAddress(w http.ResponseWriter, r *http.Request) {
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
