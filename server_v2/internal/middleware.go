package internal

import (
	"danger-dodgers/pkg/auth"

	"github.com/gofiber/fiber/v2"
)

type HTTPAuthenticator struct {
	authenticator auth.Authenticator
}

func NewHTTPAuthenticator(authenticator auth.Authenticator) *HTTPAuthenticator {
	return &HTTPAuthenticator{
		authenticator: authenticator,
	}
}

func (handler *HTTPAuthenticator) Authenticate(f func(c *fiber.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bearer, err := trimBearer(c.GetReqHeaders()["Authorization"])
		if err != nil {
			return err
		}

		payload, err := handler.authenticator.Validate(bearer)
		if err != nil {
			return err
		}

		if payload.Type != auth.Authentication {
			return &InvalidTokenTypeError{}
		}

		if payload.UserIdentifier != c.Params("id") {
			return &InvalidUserIDError{}
		}

		return f(c)
	}
}

func (handler *HTTPAuthenticator) Provide(f func(c *fiber.Ctx, id string) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bearer, err := trimBearer(c.GetReqHeaders()["Authorization"])
		if err != nil {
			return err
		}

		payload, err := handler.authenticator.Validate(bearer)
		if err != nil {
			return err
		}

		if payload.Type != auth.Authentication {
			return &InvalidTokenTypeError{}
		}

		return f(c, payload.UserIdentifier)
	}
}