package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/iots1/vertex-diagram/domain"
)

type DiagramHandler struct {
	AUsecase domain.DiagramUsecase
}

func NewDiagramHandler(app *fiber.App, us domain.DiagramUsecase) {
	handler := &DiagramHandler{
		AUsecase: us,
	}
	
	api := app.Group("/api")
	api.Get("/diagrams", handler.Fetch)
	api.Get("/diagrams/:id", handler.GetByID)
	api.Post("/diagrams", handler.Save)
	api.Delete("/diagrams/:id", handler.Delete)
}

func (h *DiagramHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.AUsecase.Delete(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func (h *DiagramHandler) Fetch(c *fiber.Ctx) error {
	list, err := h.AUsecase.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

func (h *DiagramHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := h.AUsecase.GetOne(c.Context(), id)
	if err != nil {
		log.Printf("Error fetching diagram %s: %v", id, err)
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(item)
}

func (h *DiagramHandler) Save(c *fiber.Ctx) error {
	var d domain.Diagram
	if err := c.BodyParser(&d); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.AUsecase.Save(c.Context(), &d)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}