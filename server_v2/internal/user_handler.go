package internal

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Get(c *fiber.Ctx) error {

	entity, err := handler.service.Get(c.Params("id"))
	if err != nil {
		return err
	}

	return c.JSON(entity)
}

func (handler *UserHandler) Delete(c *fiber.Ctx) error {

	user := &User{}

	err := c.BodyParser(user)
	if err != nil {
		return err
	}

	err = handler.service.Delete(user)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}

func (handler *UserHandler) Create(c *fiber.Ctx) error {

	user := &User{}

	err := c.BodyParser(user)
	if err != nil {
		return err
	}

	err = handler.service.Create(user)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}

func (handler *UserHandler) Authenticate(c *fiber.Ctx) error {
	user := &User{}

	err := c.BodyParser(user)
	if err != nil {
		return err
	}

	token, err := handler.service.Token(user)
	if err != nil {
		return err
	}

	return c.JSON(token)
}