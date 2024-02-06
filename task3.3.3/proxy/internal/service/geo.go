package service

import (
	"context"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
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
type Geocoder interface {
	SearchService() ([]*Address, error)
	GeocodeadressService(lat, lon string) (*Address, error)
}

type GeocodeService struct {
}

func (g GeocodeService) SearchService() ([]*Address, error) {
	var request SearchRequest
	cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "52a132510af242610a33fea8352874a271dbfebc",
		SecretKeyValue: "c79449d80f97afab795c7f1eda5a746d32391831"}))
	addresses, err := cleanApi.Address(context.Background(), request.Query)
	if err != nil {
		return nil, err
	}
	result := make([]*Address, len(addresses))
	for i, a := range addresses {
		result[i] = &Address{Lat: a.GeoLat, Lon: a.GeoLon}
	}
	return result, nil
}

func (g *GeocodeService) GeocodeadressService(lat, lon string) (*Address, error) {
	return &Address{Lat: lat, Lon: lon}, nil
}
