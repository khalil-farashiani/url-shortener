package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/khalil-farashiani/url-shortener/internal/drivers"
	"github.com/khalil-farashiani/url-shortener/internal/models/auth"
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"github.com/khalil-farashiani/url-shortener/internal/utils"
	"github.com/labstack/echo/v4"
)

const (
	userAssets = `assets/user/`
	message    = `Hi you get this email for reset your password
please click on the follwing link to reset your password 
if you find this email not for you igonore it
`
)

var (
	from     = utils.GetEnv("FROM", "example@gmail.com")
	password = utils.GetEnv("EMAIL_PASS", "12345678")

	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
	restErr  = make(chan error, 0)
)

func SendEmail(e auth.Email) {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, e.To, e.Message)
	if err != nil {
		fmt.Println(err.Error())
		restErr <- err
	}
	restErr <- nil
}
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

	if err := drivers.DB.First(&user, userId).Error; err != nil {
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

	if err := drivers.DB.Delete(&user, userId).Error; err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewBadRequestError("user not found"))
	}

	return c.JSON(http.StatusOK, "user deleted")
}

func UpdateUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!")
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user := &user.User{}
	if err := drivers.DB.First(&user, "username = ?", username).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Please provide valid login details"))
	}
	if utils.GetMD5(password) != user.Password {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Please provide valid login details"))
	}
	ts, err := createToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	saveErr := createAuth(user.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
	return nil
}

func ForgetPassword(c echo.Context) error {
	duration := time.Now().Add(time.Minute * 15).Unix()
	linkEx := time.Unix(duration, 0)
	now := time.Now()

	u := &user.User{}
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("invalid json body"))
	}
	if u.Email == nil {
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("email is requierd field"))
	}
	if err := drivers.DB.Where("email = ?", *u.Email).First(&u).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("user not found"))
	}
	uniqueStr := CreateUniqueLink(20)
	link := domain + uniqueStr

	setLinkErr := drivers.Client.Set(uniqueStr, strconv.Itoa(int(u.ID)), linkEx.Sub(now)).Err()
	if setLinkErr != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("we have problem to send email please try later"))
	}
	e := auth.Email{
		To:      []string{*u.Email},
		Message: []byte(message + link),
	}
	go SendEmail(e)
	err := <-restErr
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("unable to send email please try latar"))
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
}

func ResetPassword(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "implement me!!")
}
