package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"github.com/khalil-farashiani/url-shortener/internal/utils"
	"github.com/labstack/echo/v4"
)

const (
	userAssets = `assets/user/`
)

func getUserId(userIdParam string) (int64, *utils.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		err := utils.NewBadRequestError("user id should be a number")
		return 0, err
	}

	return userId, nil
}

func CreateUser(c echo.Context) error {
	// create a user struct
	var user = &user.User{}

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Get avatar
	avatar, err := c.FormFile("avatar")
	if avatar != nil && err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			utils.NewInternalServerError("unable to save avatar"),
		)
	} else if avatar != nil {
		src, err := avatar.Open()
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				utils.NewInternalServerError("unable to save avatar"),
			)
		}
		defer src.Close()

		dst, err := os.Create(userAssets + avatar.Filename)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				utils.NewInternalServerError("unable to save avatar"),
			)
		}

		if _, err = io.Copy(dst, src); err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				utils.NewInternalServerError("unable to save avatar"),
			)
		}

		defer dst.Close()
		userAvatar := userAssets + avatar.Filename
		user.Avatar = &userAvatar
	}

	user.Username = username
	user.Password = utils.GetMD5(password)
	user.Email = &email

	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError(err.Error()))
	}
	err = drivers.DB.Create(user).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("this user already exists"))
	}
	return c.JSON(http.StatusCreated, user.Marshall())
}

func GetUser(c echo.Context) error {
	idParam := c.Param("user_id")
	userId, getErr := getUserId(idParam)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}
	user := &user.User{}
	err := drivers.DB.First(&user, userId).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewBadRequestError("user not found"))
	}

	return c.JSON(http.StatusOK, user.Marshall())
}

func DeleteUser(c echo.Context) error {
	idParam := c.Param("user_id")
	userId, getErr := getUserId(idParam)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	user := &user.User{}
	err := drivers.DB.Delete(&user, userId).Error
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewBadRequestError("user not found"))
	}

	return c.JSON(http.StatusOK, "user deleted")
}

func SearchUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func UpdateUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}
