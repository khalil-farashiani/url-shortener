package main

import "github.com/khalil-farashiani/url-shortener/internal/handlers"

func routes() {
	app.GET("/ping", handlers.Ping)
	//url routes
	app.POST("/urls/", handlers.CreateUrl)
	app.GET("/urls/:url_id", handlers.GetUrl)
	app.PUT("/urls", handlers.UpdateUrl)
	app.PATCH("/urls", handlers.UpdateUrl)
	app.DELETE("/urls", handlers.DeleteUrl)
	app.GET("/urls/?created_at", handlers.SearchUrl)
	//user routes
	app.POST("/users/", handlers.CreateUser)
	app.GET("/users/:user_id", handlers.GetUser)
	app.PUT("/users", handlers.UpdateUser)
	app.PATCH("/users", handlers.UpdateUser)
	app.POST("/users/login/", handlers.Login)
	app.DELETE("/users/:user_id", handlers.DeleteUser)
}
