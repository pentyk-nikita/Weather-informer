package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/pentyk-nikita/Weather-informer/internal/domain/models"
	"github.com/pentyk-nikita/Weather-informer/pkg/config"
)

type Logger interface {
	Info(string)
	Debug(string)
	Error(string, error)
	Warn(string)
}

type WeatherInfo interface {
	GetTemperature(float64, float64) (models.TempInfo, error)
}

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, data []byte) error
}

type cliApp struct {
	l     Logger
	wi    WeatherInfo
	cache Cache
	config config.Config
}

func New(l Logger, cache Cache, wi WeatherInfo, cfg config.Config) *cliApp { 
	return &cliApp{
		l:      l,
		cache:  cache,
		wi:     wi,
		config: cfg, 
	}
}

func (c *cliApp) Run() error {
	type Current struct {
		Temp float32 `json:"temperature_2m"`
	}

	type Response struct {
		Curr Current `json:"current"`
	}

	var response Response

	cacheKey := fmt.Sprintf("weather_%.4f_%.4f", 53.6688, 23.8223)

	if cachedData, found := c.cache.Get(cacheKey); found {
		c.l.Debug("cache hit - using cached data")

		if err := json.Unmarshal(cachedData, &response); err != nil {
			c.l.Error("can't unmarshal cached data", err)
		} else {
			fmt.Printf(
				"Температура воздуха - %.2f градусов цельсия (из кеша)\n",
				response.Curr.Temp,
			)
			return nil
		}
	}

	c.l.Debug("cache miss - fetching from API")

	params := fmt.Sprintf(
		"latitude=%f&longitude=%f&current=temperature_2m",
		53.6688,
		23.8223,
	)

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s", params)
	c.l.Debug(fmt.Sprintf("url was generated success - %s", url))

	resp, err := http.Get(url)
	if err != nil {
		c.l.Error("can't get weather data", err)
		customErr := errors.New("can't get weather data from openmeteo")
		return errors.Join(customErr, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.l.Error("can't close body", err)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.l.Error("can't read data from body", err)
		customErr := errors.New("can't read data from response")
		return errors.Join(customErr, err)
	}

	c.l.Debug(fmt.Sprintf("data was readed successfully size - %d", len(data)))

	if err := c.cache.Set(cacheKey, data); err != nil {
		c.l.Warn(fmt.Sprintf("failed to save to cache: %v", err))
	}

	if err := json.Unmarshal(data, &response); err != nil {
		c.l.Error("can't unmarshal json data", err)
		customErr := errors.New("can't unmarshal data from response")
		return errors.Join(customErr, err)
	}

	fmt.Printf(
		"Температура воздуха - %.2f градусов цельсия\n",
		response.Curr.Temp,
	)
	return nil
}