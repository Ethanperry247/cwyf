package internal

import "github.com/gofiber/fiber/v2"

type DirectionsHandler struct {
	service *DirectionsService
}

func NewDirectionsHandler(service *DirectionsService) *DirectionsHandler {
	return &DirectionsHandler{
		service: service,
	}
}

func (handler *DirectionsHandler) Route(c *fiber.Ctx, id string) error {

	request := &DirectionsRequest{}

	err := c.BodyParser(request)
	if err != nil {
		return err
	}

	directions, err := handler.service.Route(request)
	if err != nil {
		return err
	}

	return c.JSON(directions)
}