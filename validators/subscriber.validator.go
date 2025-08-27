package validators

import (
	"neon/dto"
	"neon/utilities"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type SubscriberValidator struct{}

func (sv *SubscriberValidator) ValidateCreate(
	context echo.Context,
) (*dto.CreateSubscriberDto, error) {
	csDto := new(dto.CreateSubscriberDto)
	if err := context.Bind(csDto); err != nil {
		return nil, utilities.ThrowError(http.StatusBadRequest, "MALFORMED_REQUEST", err.Error())
	}

	if err := context.Validate(csDto); err != nil {
		return nil, err
	}

	csDto.Email = strings.ToLower(csDto.Email)
	return csDto, nil
}
