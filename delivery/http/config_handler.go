package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iots1/vertex-diagram/domain"
)

type ConfigHandler struct {
	ConfigUsecase domain.ConfigUsecase
}

func NewConfigHandler(app *fiber.App, uc domain.ConfigUsecase) {
	handler := &ConfigHandler{ConfigUsecase: uc}
	api := app.Group("/api")
	api.Get("/config", handler.Get)
	api.Post("/config", handler.Save)
}

func (h *ConfigHandler) Get(c *fiber.Ctx) error {
	res, err := h.ConfigUsecase.Get(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ConfigHandler) Save(c *fiber.Ctx) error {
	conf := new(domain.Config)
	if err := c.BodyParser(conf); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.ConfigUsecase.Save(c.Context(), conf); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(conf)
}
