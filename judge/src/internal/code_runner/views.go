package code_runner

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func runCodeView(c *fiber.Ctx) error {
	go func() {
		// Run the Docker container
		err := runDockerProject()
		if err != nil {
			log.Printf("Error running Docker project: %v", err)
		}
	}()
	return c.SendString("Docker project started!")
}
