package adapters

import (
	"go/hexagonal/core"

	"github.com/gofiber/fiber/v2"
)

type HttpOrderHandler struct {
	service core.OrderService
}

func NewGormHttpOrderHandler(service core.OrderService) *HttpOrderHandler {
	return &HttpOrderHandler{service: service}
}

func (h *HttpOrderHandler) CreateOrder(c *fiber.Ctx) error {
	var order core.Order
	if err := c.BodyParser(&order); err != nil {
		return err
	}

	if err := h.service.CreateOrder(order); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}
