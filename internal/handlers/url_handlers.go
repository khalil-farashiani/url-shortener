package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func CreateUrl(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func GetUrl(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func DeleteUrl(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func SearchUrl(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func UpdateUrl(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}
