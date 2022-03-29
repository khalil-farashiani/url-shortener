package main

import (
	"fmt"

	"github.com/khalil-farashiani/url-shortener/internal/models/url"
	"github.com/khalil-farashiani/url-shortener/internal/models/user"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	fmt.Println("im here")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	//migrate th User and Url entity schema
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&url.Url{})

	routes()
	app.Logger.Fatal(app.Start(portNumber))
	return nil
}
