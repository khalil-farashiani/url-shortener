package handlers

import (
	"fmt"
	"net/http"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	// Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error while getting avatar %s", err.Error()))
	}
	user := &user.User{
		Username: username,
		Email:    &email,
		Password: password,
		Avatar:   &avatar.Filename,
	}
	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	drivers.DB.Create(user)
	return c.JSON(http.StatusCreated, "user created")
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
