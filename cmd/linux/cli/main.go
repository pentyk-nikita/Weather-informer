package main

import (
	"fmt"
	"os"

	"github.com/pentyk-nikita/Weather-informer/internal/pkg/app/cli"
)

func main() {
	app := cli.New()
	err := app.Run()
	if err != nil {
		fmt.Printf("Some error - %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
