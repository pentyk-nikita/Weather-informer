package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Logger interface {
	Info(string)
	Debug(string)
	Error(string, error)
}

type cliApp struct {
	l Logger
}

func New(l Logger) *cliApp {
	return &cliApp{
		l: l,
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

	params := fmt.Sprintf(
		"latitude=%f&longitude=%f&current=temperature_2m",
		53.6688,
		23.8223,
	)

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s", params)
	c.l.Debug(fmt.Sprintf("url was generated success - %s", url))

	resp, err := http.Get(url)
	if err != nil {
		c.l.Error("can`t get weather data", err)
		customErr := errors.New("can`t get weather data from openmeteo")
		return errors.Join(customErr, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.l.Error("can`t close body", err)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.l.Error("can`t read data from body", err)
		customErr := errors.New("can`t read data from response")
		return errors.Join(customErr, err)
	}

	c.l.Debug(fmt.Sprintf("data was readed successfuly size - %d",
		len(data)))

	if err := json.Unmarshal(data, &response); err != nil {
		c.l.Error("can`t unmarshal json data", err)
		customErr := errors.New("can`t unmarshal data from response")
		return errors.Join(customErr, err)
	}

	fmt.Printf(
		"Температура воздуха - %.2f градусов цельсия\n",
		response.Curr.Temp,
	)

	return nil
}
