package middleware

import (
	"neon/utilities"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		authorization := context.Request().Header.Get("Authorization")
		authHeaderSplits := strings.Split(authorization, " ")
		authError := utilities.ThrowError(http.StatusUnauthorized, "AUTH_001", "Not authenticated")

		if len(authHeaderSplits) != 2 {
			return authError
		}

		token := authHeaderSplits[1]
		decodedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !decodedToken.Valid {
			return authError
		}
		subject, _ := decodedToken.Claims.GetSubject()

		if subject != os.Getenv("ADMIN_EMAIL") {
			return authError
		}

		return next(context)
	}
}
