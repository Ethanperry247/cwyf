package internal

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PositionHandler struct {
	service *PositionService
}

func NewPositionHandler(service *PositionService) *PositionHandler {
	return &PositionHandler{
		service: service,
	}
}

func (handler *PositionHandler) Create(c *fiber.Ctx, id string) error {
	position := &ActivityPosition{}

	err := c.BodyParser(position)
	if err != nil {
		return err
	}

	position.UserID = id

	err = handler.service.Create(position)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}

func (handler *PositionHandler) List(c *fiber.Ctx, id string) error {
	res, err := handler.service.ListByActivity(&Activity{
		ID:     c.Params("id"),
		UserID: id,
	})
	if err != nil {
		return err
	}

	return c.JSON(res)
}
