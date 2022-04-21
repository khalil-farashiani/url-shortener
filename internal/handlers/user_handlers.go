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
	"github.com/khalil-farashiani/url-shortener/logger"
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

// CreateUser godoc
// @Summary      create an user
// @Description  create an user via username and password
// @Tags         users
// @Param        username  formData      string  true  "username"
// @Param        password  formData      string  true  "password"
// @Accept       json
// @Produce      json
// @Success      200  {object}  user.User{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/ [post]
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

	if err := user.Validate("save"); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError(err.Error()))
	}
	err = drivers.DB.Create(user).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("this user already exists"))
	}
	return c.JSON(http.StatusCreated, user.Marshall())
}

// GetUser godoc
// @Summary      Show an user
// @Description  get string by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  user.User{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/{id} [get]
func GetUser(c echo.Context) error {
	idParam := c.Param("user_id")
	userId, getErr := getUserId(idParam)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}
	user := &user.User{}

	if err := drivers.DB.First(&user, userId).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewBadRequestError("user not found"))
	}

	return c.JSON(http.StatusOK, user.Marshall())
}

// DeleteUser godoc
// @Summary      delete an user
// @Description  delete an user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {string}  string "user successfully deleted"
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/{id} [delete]
func DeleteUser(c echo.Context) error {
	idParam := c.Param("user_id")
	userId, getErr := getUserId(idParam)
	if getErr != nil {
		return c.JSON(getErr.Status, getErr)
	}

	user := &user.User{}
	if err := drivers.DB.Delete(&user, userId).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewBadRequestError("user not found"))
	}

	return c.JSON(http.StatusOK, "user deleted")
}

// UpdateUser godoc
// @Summary      update an user
// @Description  update an user with PUT or PATCH method
// @Tags         users
// @Param        phone_number  formData  string  false  "phone_number"
// @Param        password  formData  string  false  "password"
// @Param        email  formData  string  false  "email"
// @Accept       json
// @Produce      json
// @Success      200  {object}  user.User{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/{id} [patch]
func UpdateUser(c echo.Context) error {
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
	userIdParam, idErr := getUserId(c.Param("user_id"))
	if uint64(userIdParam) != uint64(userId) {
		logger.Logger.Info(`user ` + strconv.Itoa(int(userIdParam)) + `try to access an endpoint that doesn't have permission`)
		return c.JSON(http.StatusForbidden, utils.NewForbiddenError("You don't have permission to do this action"))
	}
	if idErr != nil {
		logger.Logger.Info(idErr.Message)
		return c.JSON(idErr.Status, idErr)
	}

	u := user.User{}
	current := user.User{}
	current.ID = uint64(userId)
	u.ID = uint64(userId)

	if err := drivers.DB.First(&current, userId).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewBadRequestError("user not found"))
	}

	if err := c.Bind(&u); err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("invalid json body"))
	}

	isPartial := c.Request().Method == http.MethodPatch
	if isPartial {
		if u.Phonenumber != nil {
			current.Phonenumber = u.Phonenumber
		}
		if u.Password != "" {
			current.Password = utils.GetMD5(u.Password)
		}
		if u.Username != "" {
			current.Username = u.Username
		}

	} else {
		current.Phonenumber = u.Phonenumber
		current.Username = u.Username
	}

	if err := current.Validate("update"); err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError(err.Error()))
	}
	if err := drivers.DB.Model(u).Updates(user.User{Password: current.Password, Phonenumber: current.Phonenumber, Username: current.Username}).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError(err.Error()))
	}
	logger.Logger.Info(`user with id` + strconv.Itoa(int(userId)) + `updated`)
	return c.JSON(http.StatusOK, current.Marshall())
}

// Login godoc
// @Summary      login an user
// @Description  login an user
// @Tags         users
// @Param        username  formData      string  true  "username"
// @Param        password  formData      string  true  "password"
// @Accept       mpfd
// @Produce      json
// @Success      200  {object}  user.User{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/login/ [post]
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user := &user.User{}
	if err := drivers.DB.First(&user, "username = ?", username).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Please provide valid login details"))
	}
	if utils.GetMD5(password) != user.Password {
		logger.Logger.Info("unsuccessful login")
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Please provide valid login details"))
	}
	ts, err := createToken(user.ID)
	if err != nil {
		logger.Logger.Warn(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	saveErr := createAuth(user.ID, ts)
	if saveErr != nil {
		logger.Logger.Warn(err.Error())
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	logger.Logger.Info("successfully login")
	c.JSON(http.StatusOK, tokens)
	return nil
}

// Logout godoc
// @Summary      logout an user
// @Description  logout an user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  string	string "user successfully logout"
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/logout/ [post]
func Logout(c echo.Context) error {
	au, err := extractTokenMetadata(c.Request())
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	deleted, delErr := deleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		logger.Logger.Info(delErr.Error())
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	logger.Logger.Info("Successfully logged out")
	return c.JSON(http.StatusOK, "Successfully logged out")
}

// ForgetPassword godoc
// @Summary      forget password
// @Description  ForgetPassword send a token via sms or email
// @Tags         users
// @Accept       mpfd
// @Produce      json
// @Param        phone_number  formData  string  false  "phone_number"
// @Param        email  formData  string  false  "email"
// @Success      200  {object}  user.User{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/forget-password/ [post]
func ForgetPassword(c echo.Context) error {
	// Set a expire time for redis to expire token key
	duration := time.Now().Add(time.Minute * 15).Unix()
	linkEx := time.Unix(duration, 0)
	now := time.Now()

	// create a new link unique token
	uniqueStr := CreateUniqueLink(20)
	link := domain + "users/reset?token=" + uniqueStr

	u := &user.User{}
	if err := c.Bind(u); err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("invalid json body"))
	}

	// get "via" key in header request to define the method for send new password
	via := c.Request().Header.Get("via")
	if via == "" {
		logger.Logger.Info("unsuccessfully request for forget password")
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("you should pass via key header"))
	} else if via == "email" {
		if u.Email == nil {
			logger.Logger.Info("unsuccessfully request for forget password")
			return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("email is requierd field"))
		}
		// send new password via email
		// set a new key value in redis client
		setLinkErr := drivers.Client.Set(uniqueStr, strconv.Itoa(int(u.ID)), linkEx.Sub(now)).Err()
		if setLinkErr != nil {
			logger.Logger.Error(setLinkErr.Error())
			return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("we have problem to send token please try later"))
		}

		// get user with email
		if err := drivers.DB.Where("email = ?", *u.Email).First(&u).Error; err != nil {
			logger.Logger.Info(err.Error())
			return c.JSON(http.StatusNotFound, utils.NewNotFoundError("user not found"))
		}
		e := auth.Email{
			To:      []string{*u.Email},
			Message: []byte(message + link),
		}
		//send email
		go SendEmail(e)
		err := <-restErr
		if err != nil {
			logger.Logger.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("unable to send email please try latar"))
		}
		logger.Logger.Info(`an email was sent for forget password to user`)
		return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
	} else if via == "sms" {
		// Send password via sms
		if u.Phonenumber == nil {
			logger.Logger.Info("unsuccessfully request for forget password")
			return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("phonenumber is requierd field"))
		}
		setLinkErr := drivers.Client.Set(uniqueStr, strconv.Itoa(int(u.ID)), linkEx.Sub(now)).Err()
		if setLinkErr != nil {
			logger.Logger.Error(setLinkErr.Error())
			return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("we have problem to send token please try later"))
		}

		// get user with phonenumber
		if err := drivers.DB.Where("phonenumber = ?", u.Phonenumber).First(&u).Error; err != nil {
			logger.Logger.Info(err.Error())
			return c.JSON(http.StatusNotFound, utils.NewNotFoundError("user not found"))
		}
		if err := utils.SendSms(message+link, *u.Phonenumber); err != nil {
			logger.Logger.Error(err.Error)
			return c.JSON(http.StatusInternalServerError, utils.NewInternalServerError("unable to send sms, plaese try later"))
		}
		logger.Logger.Info("an SMS was sent for forget password to user")
		return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
	}
	return c.JSON(http.StatusBadRequest, utils.NewBadRequestError("invalid value in via key header"))
}

// ResetPassword godoc
// @Summary      reset password
// @Description  ResetPassword validate token that user click on it and send a new password
// @Tags         users
// @Accept       mpfd
// @Produce      json
// @Success      200  {string}  string "a new password created"
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/reset/ [get]
func ResetPassword(c echo.Context) error {
	token := c.QueryParam("token")
	userId, err := drivers.Client.Get(token).Result()
	if err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("invalid token param"))
	}

	newPass := CreateUniqueLink(9)
	if err := drivers.DB.Model(&user.User{}).Where("id = ?", userId).Update("password", utils.GetMD5(newPass)).Error; err != nil {
		logger.Logger.Info(err.Error())
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("user not found"))
	}
	// delete token before show the response
	pipe := drivers.Client.Pipeline()
	pipe.Del(token)
	pipe.Exec()

	return c.JSON(http.StatusOK, map[string]string{
		"message":      "OK this is your new password you can change it in your profile",
		"new_password": newPass,
	})
}

// EnableSpecialUser godoc
// @Summary      enable premium
// @Description  EnableSpecialUser enable some features for users like create shorter link
// @Tags         users
// @Accept       mpfd
// @Produce      json
// @Success      200  {object}  user.User{}
// @Failure      400  {object}  utils.RestErr{}
// @Failure      401  {object}  utils.RestErr{}
// @Failure      404  {object}  utils.RestErr{}
// @Failure      500  {object}  utils.RestErr{}
// @Router       /users/premium/ [get]
func EnableSpecialUser(c echo.Context) error {
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
	user := user.User{ID: userId}
	drivers.DB.Model(&user).Update("is_special", true)
	drivers.DB.First(&user)
	// call zarinpall method to get the response
	// if response show not successful
	// return utils.NewBadRequest("payment is not successful")
	//else
	logger.Logger.Info("an user successfully enable premium plan")
	return c.JSON(http.StatusOK, user.Marshall())
}
