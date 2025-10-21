package controllers

import (
	"fmt"
	"neon/utilities"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type oauthController struct {
}

func RegisterOauthRoutes(group *echo.Group) {
	oc := &oauthController{}
	group.GET("/initiate", oc.redirect)
	group.GET("/authorize", oc.callbackHandler)
}

func (oauthController *oauthController) redirect(context echo.Context) error {
	query := context.Request().URL.Query()
	query.Add("provider", "google")
	context.Request().URL.RawQuery = query.Encode()

	request := context.Request()
	response := context.Response().Writer

	gothic.Store = utilities.GetOauthSessionStore()

	if gothUser, err := gothic.CompleteUserAuth(response, request); err == nil {
		return context.JSON(http.StatusOK, gothUser)
	}

	gothic.BeginAuthHandler(response, request)
	return nil
}

func (oauthController *oauthController) callbackHandler(context echo.Context) error {
	query := context.Request().URL.Query()
	query.Add("provider", "google")
	context.Request().URL.RawQuery = query.Encode()

	request := context.Request()
	response := context.Response().Writer

	gothic.Store = utilities.GetOauthSessionStore()
	user, err := gothic.CompleteUserAuth(response, request)

	if err != nil {
		fmt.Println(err)
		return utilities.ThrowError(
			http.StatusInternalServerError,
			"AUTH_003",
			"There was a problem retrieving user information from google",
		)
	}

	if strings.ToLower(user.Email) != os.Getenv("ADMIN_EMAIL") {
		return utilities.ThrowError(
			http.StatusBadRequest,
			"AUTH_001",
			"Login attempt made by an unauthorized user",
		)
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": os.Getenv("ADMIN_EMAIL"),
		"iss": os.Getenv("APP_NAME"),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString(secretKey)

	if err != nil {
		return utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"There was a problem generating the auth token",
		)
	}

	tokenCookie := new(http.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = token
	tokenCookie.Expires = time.Now().Add(time.Hour * 72)
	tokenCookie.Path = "/"
	http.SetCookie(context.Response().Writer, tokenCookie)

	return context.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
}
