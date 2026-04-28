package gui

import (
	"fmt"
	"time"

	guisettings "github.com/pentyk-nikita/Weather-informer/internal/domain/gui_settings"
	"github.com/pentyk-nikita/Weather-informer/pkg/config"
)

type Logger interface {
	Info(string)
	Debug(string)
	Error(string, error)
}

type WeatherInfo interface {
	GetTemperature(float64, float64) (float32, error)
}

type GUIApp struct {
	l          Logger
	p          guisettings.Provider
	wi         WeatherInfo
	conf       config.Config
	window     guisettings.Window
	textWidget guisettings.TextWidget
}

func New(l Logger, p guisettings.Provider, wi WeatherInfo, c config.Config) *GUIApp {
	return &GUIApp{
		l:    l,
		p:    p,
		wi:   wi,
		conf: c,
	}
}

func (g *GUIApp) Run() error {
	size := guisettings.NewWS(400, 300)
	window, err := g.p.CreateWindow("Weather Informer", size)
	if err != nil {
		g.l.Error("Failed to create window", err)
		return err
	}
	g.window = window

	g.textWidget = g.p.GetTextWidget("Loading weather data...")

	if err := window.SetTemperatureWidget(g.textWidget); err != nil {
		g.l.Error("Failed to set widget", err)
		return err
	}

	go func() {
		for {
			temp, err := g.wi.GetTemperature(g.conf.L.Lat, g.conf.L.Long)
			if err != nil {
				g.l.Error("Failed to get temperature", err)
				g.textWidget.SetText("Error getting weather data")
			} else {
				g.textWidget.SetText(fmt.Sprintf("Температура воздуха: %.2f °C", temp))
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	g.p.GetAppRunner().Run()
	return nil
}
