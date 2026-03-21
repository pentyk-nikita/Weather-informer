package main

import (
	"os"

	"github.com/pentyk-nikita/weather_info/internal/pkg/app/cli"
	"github.com/pentyk-nikita/weather_info/pkg/logger"
)

func main() {
	l := logger.New()
	app := cli.New(l)

	err := app.Run()
	if err != nil {
		l.Error("Some error", err)
		os.Exit(1)
	}

	os.Exit(0)
}
