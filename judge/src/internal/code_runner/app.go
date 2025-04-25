package code_runner

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func StartListening(port int) {

	app := fiber.New(fiber.Config{})
	// Landing
	// app.Get("/", landingView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
