package main

import (
	"os"
	"time"

	"github.com/pentyk-nikita/Weather-informer/internal/pkg/app/cli"
	"github.com/pentyk-nikita/Weather-informer/pkg/cache"
	"github.com/pentyk-nikita/Weather-informer/pkg/logger"
)

func main() {
	l := logger.New()

	c, err := cache.New("./cache", 10*time.Minute)
	if err != nil {
		l.Error("Failed to create cache", err)
		os.Exit(1)
	}

	app := cli.New(l, c)

	err = app.Run()
	if err != nil {
		l.Error("Some error", err)
		os.Exit(1)
	}

	os.Exit(0)
}