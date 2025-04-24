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
	app.Post("/signin", handleSigninView)
	app.Post("/signup", handleSignupView)
	app.Get("/signout", signoutView)

	// problemset
	app.Get("/problemset", isAuthenticated, problemsetView)
	app.Get("/problemset/add", isAuthenticated, addProblemView)
	app.Post("/problemset/add", isAuthenticated, handleAddProblemView)

	app.Get("/problemset/:id", isAuthenticated, showProblemView)
	app.Get("/problemset/:id/dl/:filename", isAuthenticated, downloadProblemInOutFiles)
	app.Get("/problemset/:id/:command", isAuthenticated, handlePublishProblemView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
