package main

import "github.com/prometheus/client_golang/prometheus"

// Создание метрик для времени выполнения запросов к эндпоинтам
var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	// Создание метрик для количества запросов к эндпоинтам
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"endpoint"},
	)

	// Создание метрик для времени обращения в кэш
	cacheDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_duration_seconds",
			Help:    "Duration of cache requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Создание метрик для времени обращения в базу данных
	dbDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_duration_seconds",
			Help:    "Duration of database requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Создание метрик для времени обращения во внешний API
	externalAPIDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_duration_seconds",
			Help:    "Duration of external API requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

// Регистрация метрик
func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(cacheDuration)
	prometheus.MustRegister(dbDuration)
	prometheus.MustRegister(externalAPIDuration)
}
