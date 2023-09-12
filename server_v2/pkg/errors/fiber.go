package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var FiberMappings = map[Group]int{
	BAD_REQUEST: http.StatusBadRequest,
	NOT_FOUND:   http.StatusNotFound,
	FORBIDDEN:   http.StatusForbidden,
}

type FiberErrorHandler struct{}

func New() *FiberErrorHandler {
	return &FiberErrorHandler{}
}

func (handler *FiberErrorHandler) HandleError(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	internal, ok := err.(Error)
	if ok {
		mapping, ok := FiberMappings[internal.Code()]
		if ok {
			code = mapping
		}
	}

	if code == fiber.StatusInternalServerError {
		fmt.Println("internal error")
	} else {
		fmt.Println("non-internal error")
	}


	err = c.Status(code).Send([]byte(err.Error()))
	if err != nil {
		fmt.Println("could not send error")
	}

	return nil
}