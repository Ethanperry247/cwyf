package internal

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ActivityHandler struct {
	service *ActivityService
}

func NewActivityHandler(service *ActivityService) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

func (handler *ActivityHandler) Delete(c *fiber.Ctx, id string) error {
	err := handler.service.Delete(&Activity{
		ID:     c.Params("id"),
		UserID: id,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}

func (handler *ActivityHandler) Create(c *fiber.Ctx, id string) error {
	activity := &Activity{}

	err := c.BodyParser(activity)
	if err != nil {
		return err
	}

	activity.UserID = id

	err = handler.service.Create(activity)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}

func (handler *ActivityHandler) List(c *fiber.Ctx, id string) error {
	res, err := handler.service.ListByUser(id)
	if err != nil {
		return err
	}

	return c.JSON(res)
}