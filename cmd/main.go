package main

import (
	"fmt"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

const (
	portNumber = ":8080"
)

var app = echo.New()

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
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

	routes()

	app.Logger.Fatal(app.Start(portNumber))
	return nil
}
