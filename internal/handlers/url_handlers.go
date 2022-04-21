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
	"github.com/khalil-farashiani/url-shortener/logger"
	"github.com/labstack/echo/v4"
)

const (
	domain = "127.0.0.1:8080/"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func CreateUniqueLink(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// CreateUrl godoc
// @Summary      create an short url
// @Description  CreateUrl create a short url
// @Tags         urls
// @Param        source  formData  string  true  "source"
// @Accept       json
// @Produce      json
// @Success      200  {object}  url.Url{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /urls/ [post]
func CreateUrl(c echo.Context) error {
	url := &url.Url{}
	source := c.FormValue("source")
	tokenAuth, err := extractTokenMetadata(c.Request())
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}
	url.Source = source
	user := &user.User{}

	if err := drivers.DB.First(&user, userId).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("user not found"))
	}
	url.User = *user
	url.UserID = userId
	if user.IsSpecial == false {
		url.ShortUrl = domain + CreateUniqueLink(7)
	} else {
		url.ShortUrl = domain + CreateUniqueLink(5)
	}

	updateErr := drivers.DB.Create(url).Error
	if updateErr != nil {
		logger.Logger.Info(updateErr.Error())
		return c.JSON(http.StatusNotFound, utils.NewInternalServerError("have an issue in create shortlink"))
	}
	logger.Logger.Info("new short url was created")
	return c.JSON(http.StatusCreated, url.Marshall())
}

// GetUrl godoc
// @Summary      get url
// @Description  GetUrl get the main url to redirect
// @Tags         urls
// @Param        url  path  string true  "url"
// @Accept       json
// @Produce      json
// @Success      200  {object}  url.Url{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /urls/{url} [get]
func GetUrl(c echo.Context) error {
	url := &url.Url{}
	urlParam := c.Param("url")
	if err := drivers.DB.First(&url, "short_url = ?", domain+urlParam).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("url not found"))
	}
	return c.JSON(http.StatusOK, map[string]string{"result": url.Source})
}

// DeleteUrl godoc
// @Summary      delete an url
// @Description  delete an url with
// @Tags         urls
// @Accept       json
// @Produce      json
// @Success      200  {object}  url.Url{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /urls/{url} [delete]
func DeleteUrl(c echo.Context) error {
	urlParam := c.Param("url")
	tokenAuth, err := extractTokenMetadata(c.Request())
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	if err := drivers.DB.Where(map[string]interface{}{"short_url": domain + urlParam, "user_id": userId}).Delete(&url.Url{}).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("url not found"))
	}
	logger.Logger.Info("an url was deleted")
	return c.JSON(http.StatusOK, map[string]string{"message": "deleted"})
}

// Myurls godoc
// @Summary      show list of urls
// @Description  MyUrls return a list of user's urls
// @Tags         urls
// @Accept       json
// @Produce      json
// @Success      200  {array}  url.Url{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /urls/my-links [get]
func MyUrls(c echo.Context) error {
	urls := &url.Urls{}
	tokenAuth, err := extractTokenMetadata(c.Request())
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("unauthorized"))
	}

	if err := drivers.DB.Where(&url.Url{UserID: uint64(userId)}).Find(urls).Error; err != nil {
		logger.Logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("we have problem to return the urls"))
	}
	return c.JSON(http.StatusOK, urls.Marshall())
}
