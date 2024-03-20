package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CacheDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cache_request_duration_seconds",
			Help:    "Histogram of the cache request duration in seconds.",
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
		},
	)
	CacheCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_requests_total",
			Help: "Total number of cache requests.",
		},
	)
	DbDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "db_request_duration_seconds",
			Help:    "Histogram of the database request duration in seconds.",
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
		},
	)
	DbCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_requests_total",
			Help: "Total number of database requests.",
		},
	)
	ApiDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Histogram of the external API request duration in seconds.",
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
		},
	)
	ApiCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of external API requests.",
		},
	)
)

func init() {
	prometheus.MustRegister(CacheDuration)
	prometheus.MustRegister(CacheCount)
	prometheus.MustRegister(DbDuration)
	prometheus.MustRegister(DbCount)
	prometheus.MustRegister(ApiDuration)
	prometheus.MustRegister(ApiCount)
}
