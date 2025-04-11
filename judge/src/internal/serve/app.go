package serve

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func StartListening(port int) {
	connectDatabase()
	initSessionStore()

	app := fiber.New(fiber.Config{
		Views: html.New("static/views/", ".html"),
	})

	// Landing
	app.Get("/", landingView)

	// signin
	app.Get("/signin", signinView)
	app.Post("/signin", signinSubmitView)
	app.Post("/signup", signupSubmitView)
	app.Get("/signout", signoutView)

	// problemset
	app.Get("/problemset", isAuthenticated, problemsetView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
