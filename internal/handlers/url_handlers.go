package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
	"github.com/khalil-farashiani/url-shortener/internal/models/url"
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"github.com/khalil-farashiani/url-shortener/internal/utils"
	"github.com/labstack/echo/v4"
)

const (
	domain = "127.0.0.1:8080/"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func createShortLink(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func CreateUrl(c echo.Context) error {
	url := &url.Url{}
	source := c.FormValue("source")
	tokenAuth, err := extractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}
	url.Source = source
	url.ShortUrl = domain + createShortLink(7)
	user := &user.User{}

	if err := drivers.DB.First(&user, userId).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("user not found"))
	}
	url.User = *user
	url.UserID = userId

	updateErr := drivers.DB.Create(url).Error
	if updateErr != nil {
		return c.JSON(http.StatusNotFound, utils.NewInternalServerError("have an issue in create shortlink"))
	}
	return c.JSON(http.StatusCreated, url.Marshall())
}

func GetUrl(c echo.Context) error {
	url := &url.Url{}
	urlParam := c.Param("url")
	if err := drivers.DB.First(&url, "short_url = ?", domain+urlParam).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("url not found"))
	}
	return c.JSON(http.StatusOK, map[string]string{"result": url.Source})
}

func DeleteUrl(c echo.Context) error {
	urlParam := c.Param("url")
	tokenAuth, err := extractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	if err := drivers.DB.Where(map[string]interface{}{"short_url": domain + urlParam, "user_id": userId}).Delete(&url.Url{}).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("url not found"))
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "deleted"})
}

func MyUrls(c echo.Context) error {
	urls := &url.Urls{}
	tokenAuth, err := extractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	if err := drivers.DB.Where(&url.Url{UserID: uint64(userId)}).Find(urls).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("we have problem to return the urls"))
	}
	return c.JSON(http.StatusOK, urls.Marshall())
}
