package main

import (
	"os"
	"time"

	"github.com/pentyk-nikita/Weather-informer/internal/adapters/weather"
	"github.com/pentyk-nikita/Weather-informer/internal/pkg/app/cli"
	"github.com/pentyk-nikita/Weather-informer/internal/pkg/flags"
	"github.com/pentyk-nikita/Weather-informer/pkg/cache"
	"github.com/pentyk-nikita/Weather-informer/pkg/config"
	"github.com/pentyk-nikita/Weather-informer/pkg/logger"

	pogodaby "github.com/pentyk-nikita/Weather-informer/internal/adapters/pogoda_by"
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
	
	cacheInst, err := cache.New("./cache", 10*time.Minute)
	if err != nil {
		l.Error("Failed to create cache", err)
		os.Exit(1)
	}

	wi := getProvider(c, l)
	app := cli.New(l, cacheInst, wi, c)

	err = app.Run()
	if err != nil {
		l.Error("Some error", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func getProvider(c config.Config, l cli.Logger) cli.WeatherInfo {
	var wi cli.WeatherInfo
	switch c.P.Type {
	case "open-meteo":
		wi = weather.New(l)
	case "pogoda":
		wi = pogodaby.New(l)
	default:
		wi = weather.New(l)
	}
	return wi
}