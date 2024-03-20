package service

import (
	"context"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/client"
	"hugo/task4.2.4/proxy/internal/repository"
	"hugo/task4.2.4/proxy/metrics"
	"time"
)

type Geocoder interface {
	SearchService(request repository.SearchRequest) ([]*repository.Address, error)
	GeocodeAddressService(lat, lon string) (*string, error)
}

type GeocodeService struct {
	repo repository.GeocoderRepository
}

func NewGeoSerive(repo repository.GeocoderRepository) *GeocodeService {
	return &GeocodeService{repo: repo}
}

func (g *GeocodeService) SearchService(request repository.SearchRequest) ([]*repository.Address, error) {
	historyID, err := g.repo.SaveSearchHistory(request.Query)
	if err != nil {
		return nil, err
	}

	similarAddresses, err := g.repo.Search(request.Query)
	if err != nil {
		return nil, err
	}

	if len(similarAddresses) > 0 {
		addresses := make([]*repository.Address, len(similarAddresses))
		for i, addr := range similarAddresses {
			addresses[i] = &addr
		}
		return addresses, nil
	}
	start := time.Now()
	metrics.ApiCount.Inc()

	cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(&client.Credentials{
		ApiKeyValue:    "11cb4969967b7e68ab87b57258372aefec0eb6ac",
		SecretKeyValue: "3461265109aaa28b20523e1b4dfb4d36e475fc9f"}))
	addresses, err := cleanApi.Address(context.Background(), request.Query)
	if err != nil {
		return nil, err
	}

	result := make([]*repository.Address, len(addresses))
	for i, a := range addresses {
		result[i] = &repository.Address{Lat: a.GeoLat, Lon: a.GeoLon}
		addressID, err := g.repo.SaveAddress(request.Query, a.GeoLat, a.GeoLon)
		g.repo.SaveHistorySearchAddress(historyID, addressID)
		if err != nil {
			return nil, err
		}
	}

	metrics.ApiDuration.Observe(time.Since(start).Seconds())
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
