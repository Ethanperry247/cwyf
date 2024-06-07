package internal

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service *ReportService
}

func NewReportHandler(service *ReportService) *ReportHandler {
	return &ReportHandler{
		service: service,
	}
}

func (handler *ReportHandler) Get(c *fiber.Ctx, id string) error {

	entity, err := handler.service.Get(&Report{
		ID:     c.Params("id"),
		UserID: id,
	})
	if err != nil {
		return err
	}

	return c.JSON(entity)
}

func (handler *ReportHandler) Delete(c *fiber.Ctx, id string) error {
	err := handler.service.Delete(&Report{
		ID:     c.Params("id"),
		UserID: id,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}

func (handler *ReportHandler) Create(c *fiber.Ctx, id string) error {
	report := &Report{}

	err := c.BodyParser(report)
	if err != nil {
		return err
	}

	report.UserID = id

	err = handler.service.Create(report)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}

func (handler *ReportHandler) List(c *fiber.Ctx, id string) error {
	res, err := handler.service.ListByUser(id)
	if err != nil {
		return err
	}

	return c.JSON(res)
}
