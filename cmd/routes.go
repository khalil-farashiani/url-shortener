package main

import "github.com/khalil-farashiani/url-shortener/internal/handlers"

func routes() {
	app.GET("/ping", handlers.Ping)
	//url routes
	app.POST("/urls/", handlers.CreateUrl)
	app.GET("/:url", handlers.GetUrl)
	app.DELETE("/:url", handlers.DeleteUrl)
	app.GET("/users/my-links", handlers.MyUrls)
	//user routes
	app.POST("/users/", handlers.CreateUser)
	app.GET("/users/:user_id", handlers.GetUser)
	app.PUT("/users", handlers.UpdateUser)
	app.PATCH("/users", handlers.UpdateUser)
	app.POST("/users/login/", handlers.Login)
	app.DELETE("/users/:user_id", handlers.DeleteUser)
}
