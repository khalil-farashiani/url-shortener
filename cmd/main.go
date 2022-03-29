package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	portNumber = ":8080"
)

var app = echo.New()

func main() {
	err := run()
	if err != nil {
		panic(fmt.Sprintf("Error while trying to start application: %s", err.Error()))
	}
}

func run() error {
	routes()
	app.Logger.Fatal(app.Start(portNumber))
	return nil
}
