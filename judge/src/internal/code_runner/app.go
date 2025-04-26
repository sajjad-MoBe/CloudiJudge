package code_runner

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

var queueManager *QueueManager

func StartListening(port int) {
	queueManager = NewQueueManager(20)
	app := fiber.New(fiber.Config{})

	app.Get("/run", runCodeView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
