package controllers

import (
	"neon/utilities"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type OauthController struct {
}

func (oauthController *OauthController) Redirect(context echo.Context) error {
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
