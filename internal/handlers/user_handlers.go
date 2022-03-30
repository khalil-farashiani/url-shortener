package handlers

import (
	"net/http"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	drivers.DB.Create(&user.User{})
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func GetUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func DeleteUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func SearchUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func UpdateUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}
