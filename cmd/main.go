package main

import (
	"fmt"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
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
	//connect to database
	err := drivers.ConnectSQL()
	if err != nil {
		return err
	}

	app.Static("/assets", "assets")
	// app.Use(middleware.JWT([]byte("secret")))

	routes()

	app.Logger.Fatal(app.Start(portNumber))
	return nil
}
