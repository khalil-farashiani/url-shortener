package main

import "github.com/khalil-farashiani/url-shortener/internal/handlers"

func routes() {
	app.GET("/ping", handlers.Ping)
}
