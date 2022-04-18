package main

import (
	"github.com/khalil-farashiani/url-shortener/internal/handlers"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func routes() {
	url := echoSwagger.URL("http://127.0.0.1:8080/swagger/doc.json")
	app.GET("/ping", handlers.Ping)
	app.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	//url routes
	app.POST("/urls/", handlers.CreateUrl)
	app.GET("/:url", handlers.GetUrl)
	app.DELETE("/:url", handlers.DeleteUrl)
	app.GET("/users/my-links", handlers.MyUrls)
	//user routes
	app.POST("/users/", handlers.CreateUser)
	app.GET("/users/:user_id", handlers.GetUser)
	app.PUT("/users/:user_id", handlers.UpdateUser)
	app.PATCH("/users/:user_id", handlers.UpdateUser)
	app.POST("/users/login/", handlers.Login)
	app.DELETE("/users/:user_id", handlers.DeleteUser)
	app.POST("/users/forget-password", handlers.ForgetPassword)
	app.GET("/users/reset", handlers.ResetPassword)
	app.GET("/users/premium", handlers.EnableSpecialUser)
	app.POST("/users/logout", handlers.Logout)
}
