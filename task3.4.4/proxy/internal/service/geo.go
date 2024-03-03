package service

import (
	"context"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"hugoproxy-main/task3.3.3/proxy/internal/repository"
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
type Geocoder interface {
	SearchService(request SearchRequest) ([]*Address, error)
	GeocodeAddressService(lat, lon string) (*string, error)
}

type GeocodeService struct {
	repo repository.GeocoderRepository
}

func NewGeoSerive(repo repository.GeocoderRepository) *GeocodeService {
	return &GeocodeService{repo: repo}
}

func (g *GeocodeService) SearchService(request SearchRequest) ([]*Address, error) {
	similarAddresses, err := g.repo.Search(request.Query)
	if err != nil {
		return nil, err
	}

	if len(similarAddresses) > 0 {
		addresses := make([]*Address, len(similarAddresses))
		for i, addr := range similarAddresses {
			addresses[i] = &Address{AddressText: addr}
		}
		return addresses, nil
	}

	cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "11cb4969967b7e68ab87b57258372aefec0eb6ac",
		SecretKeyValue: "3461265109aaa28b20523e1b4dfb4d36e475fc9f"}))
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

func (g *GeocodeService) GeocodeAddressService(lat, lon string) (*string, error) {
	cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "11cb4969967b7e68ab87b57258372aefec0eb6ac",
		SecretKeyValue: "3461265109aaa28b20523e1b4dfb4d36e475fc9f"}))
	addresses, err := cleanApi.Address(context.Background(), lat, lon)
	if err != nil {
		return nil, err
	}
	return &addresses[0].City, nil
}
