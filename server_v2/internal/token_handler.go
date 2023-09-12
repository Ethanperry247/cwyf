package internal

import (
	"danger-dodgers/pkg/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokenHandler struct {
	service *TokenService
	auth    auth.Authenticator
}

func NewTokenHandler(service *TokenService,
	auth auth.Authenticator) *TokenHandler {
	return &TokenHandler{
		service: service,
		auth:    auth,
	}
}

func trimBearer(bearer string) (string, error) {
	if bearer == "" {
		return "", &MissingBearerTokenError{}
	}

	components := strings.Split(bearer, " ")
	if len(components) != 2 || components[0] != "Bearer" {
		return "", &ImproperlyFormattedBearerTokenError{}
	}

	return components[1], nil
}

func (handler *TokenHandler) Create(c *fiber.Ctx) error {
	bearer, err := trimBearer(c.GetReqHeaders()["Authorization"])
	if err != nil {
		return err
	}

	payload, err := handler.auth.Validate(bearer)
	if err != nil {
		return err
	}

	if payload.Type != auth.Refresh {
		return &InvalidTokenTypeError{}
	}

	token, err := handler.service.Token(&User{
		Username: payload.UserIdentifier,
	})
	if err != nil {
		return err
	}

	return c.JSON(token)
}
