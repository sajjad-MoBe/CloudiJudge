package code_runner

import (
	"github.com/gofiber/fiber/v2"
)

func runCodeView(c *fiber.Ctx) error {
	var data Run

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}
	queueManager.Enqueue(data)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok": true,
	})
}
