package code_runner

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var queueManager *QueueManager

func StartListening(port int) {
	maxConcurrent, err := strconv.Atoi(os.Getenv("MAX_CONCURRENT_RUNS"))
	if err != nil {
		fmt.Println("invalid MAX_CONCURRENT_RUNS value, use 10 as default value")
		maxConcurrent = 10
	}
	queueManager = NewQueueManager(maxConcurrent)
	app := fiber.New(fiber.Config{})

	app.Post("/run", runCodeView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
