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
	l      Logger
	wi     WeatherInfo
	cache  Cache
	conf   config.Config
}

func New(l Logger, cache Cache, wi WeatherInfo, c config.Config) *cliApp {
	return &cliApp{
		l:     l,
		cache: cache,
		wi:    wi,
		conf:  c,
	}
}

func (c *cliApp) Run() error {
	lat := c.conf.L.Lat
	long := c.conf.L.Long

	cacheKey := fmt.Sprintf("weather_%.4f_%.4f", lat, long)

	if cachedData, found := c.cache.Get(cacheKey); found {
		c.l.Debug("cache hit - using cached data")

		var response struct {
			Curr struct {
				Temp float32 `json:"temperature_2m"`
			} `json:"current"`
		}

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

	tempInfo, err := c.wi.GetTemperature(lat, long)
	if err != nil {
		c.l.Error("can't get temperature from provider", err)
		return err
	}

	cacheData := struct {
		Current struct {
			Temperature2m float32 `json:"temperature_2m"`
		} `json:"current"`
	}{}
	cacheData.Current.Temperature2m = tempInfo.Temp
	
	jsonData, _ := json.Marshal(cacheData)
	if err := c.cache.Set(cacheKey, jsonData); err != nil {
		c.l.Warn(fmt.Sprintf("failed to save to cache: %v", err))
	}

	fmt.Printf(
		"Температура воздуха - %.2f градусов цельсия\n",
		tempInfo.Temp,
	)
	return nil
}