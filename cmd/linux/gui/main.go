package main

import (
	"os"

	"github.com/pentyk-nikita/Weather-informer/internal/adapters/weather"
	"github.com/pentyk-nikita/Weather-informer/internal/pkg/app/gui"
	"github.com/pentyk-nikita/Weather-informer/internal/pkg/flags"
	"github.com/pentyk-nikita/Weather-informer/pkg/config"
	"github.com/pentyk-nikita/Weather-informer/pkg/logger"
	
	guifyne "github.com/pentyk-nikita/Weather-informer/internal/pkg/gui/fyne"
)

func main() {
	arguments := flags.Parse()

	r, err := os.Open(arguments.Path)
	if err != nil {
		panic(err)
	}

	c, err := config.Parse(r)
	if err != nil {
		panic(err)
	}

	l := logger.New()

	wi := getProvider(c, l)

	p := guifyne.NewP()

	g := gui.New(l, p, wi, c)

	err = g.Run()
	if err != nil {
		panic(err)
	}
}

func getProvider(c config.Config, l logger.Logger) gui.WeatherInfo {
	var wi gui.WeatherInfo
	switch c.P.Type {
	case "open-meteo":
		wi = weather.New(l)
	default:
		wi = weather.New(l)
	}
	return wi
}
