package pogodaby

import (
	"encoding/json"
	"net/http"

	"github.com/pentyk-nikita/Weather-informer/internal/domain/models"
)

const url = "https://pogoda.by/api/v2/weather-fact?station=26820"

type resp struct {
	Temp float32 `json:"t"`
}

type Logger interface {
	Info(string)
	Debug(string)
	Error(string, error)
}

type pogoda struct {
	l Logger
}

func New(l Logger) *pogoda {
	return &pogoda{l: l}
}

func (p *pogoda) GetTemperature(lat, long float64) (models.TempInfo, error) {
	response, err := http.Get(url)
	if err != nil {
		p.l.Error("can't get data from pogoda.by", err)
		return models.TempInfo{}, err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			p.l.Error("can't close response body", err)
		}
	}()

	var r resp

	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		p.l.Error("can't decode JSON", err)
		return models.TempInfo{}, err
	}
	return models.TempInfo{Temp: r.Temp}, nil
}