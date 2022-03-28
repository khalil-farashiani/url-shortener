package main

import (
	"github.com/labstack/echo/v4"
)

const (
	portNumber = ":8080"
)

var app = echo.New()

func main() {
	run()
}

func run() {
	routes()
	app.Logger.Fatal(app.Start(portNumber))
}
