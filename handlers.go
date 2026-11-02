package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const cacheTTL = 12 * time.Hour

type WeatherData struct {
	Address           string `json:"address"`
	CurrentConditions struct {
		Temp float64 `json:"temp"`
	} `json:"currentConditions"`
}

type WeatherClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewWeatherClient(apiKey string) *WeatherClient {
	return &WeatherClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (c *WeatherClient) GetWeather(ctx context.Context, city string) ([]byte, error) {
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("include", "current") // не тянем прогноз на 15 дней
	q.Set("elements", "temp")   // только нужные поля
	u := "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/" +
		url.PathEscape(city) + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

type Server struct {
	Rdb    *redis.Client
	Client *WeatherClient
	group  singleflight.Group
}

func (s *Server) handlerWeather(w http.ResponseWriter, r *http.Request) {
	city := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("city")))
	if city == "" {
		http.Error(w, "city parameter is required", http.StatusBadRequest)
		return
	}
	key := "weather:" + city
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	cached, err := s.Rdb.Get(ctx, key).Bytes()
	if err == nil {
		w.Write(cached)
		return
	}
	if !errors.Is(err, redis.Nil) {
		log.Printf("redis get: %v", err)
	}

	// одинаковые параллельные запросы схлопываются в один вызов API
	v, err, _ := s.group.Do(key, func() (any, error) {
		fetchCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := s.Client.GetWeather(fetchCtx, city)
		if err != nil {
			return nil, err
		}
		if err := s.Rdb.Set(fetchCtx, key, data, cacheTTL).Err(); err != nil {
			log.Printf("redis set: %v", err)
		}
		return data, nil
	})
	if err != nil {
		log.Printf("weather fetch: %v", err)
		http.Error(w, "Failed to fetch weather", http.StatusBadGateway)
		return
	}

	w.Write(v.([]byte))
}
