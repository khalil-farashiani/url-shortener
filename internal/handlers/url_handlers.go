package handlers

import (
	"fmt"
	"net/http"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
	"github.com/khalil-farashiani/url-shortener/internal/models/url"
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"github.com/khalil-farashiani/url-shortener/internal/utils"
	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func CreateUrl(c echo.Context) error {
	url := &url.Url{}
	user := &user.User{}
	source := c.FormValue("source")
	tokenAuth, err := extractTokenMetadata(c.Request())
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	userId, err := FetchAuth(tokenAuth)
	fmt.Println(user.ID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	user.ID = userId
	url.Source = source

	updateErr := drivers.DB.Model(&user).Where("id=?", userId).Association("url").Append(&url).Error
	if updateErr != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("user not found"))
	}
	return c.JSON(http.StatusCreated, url)
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
