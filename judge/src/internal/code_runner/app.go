package code_runner

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func StartListening(port int) {

	app := fiber.New(fiber.Config{})
	// Landing
	app.Get("/run", runCodeView)

	// log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
	log.Fatal(runDockerProject())

}
