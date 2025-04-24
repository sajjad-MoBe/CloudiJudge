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

	htmlEngine := html.New("static/views/", ".html")
	htmlEngine.AddFunc("add", Add)
	htmlEngine.AddFunc("sub", Sub)
	htmlEngine.AddFunc("truncate", Truncate)

	app := fiber.New(fiber.Config{
		Views: htmlEngine,
	})

	// Landing
	app.Get("/", landingView)

	// signin
	app.Get("/signin", signinView)
	app.Post("/signin", handleSigninView)
	app.Post("/signup", handleSignupView)
	app.Get("/signout", signoutView)

	app.Get("/user/:id", isAuthenticated, showProfileView)
	app.Post("/user/:id/promote", isAuthenticated, promoteUserView)
	app.Post("/user/:id/demote", isAuthenticated, demoteUserView)

	// problemset
	app.Get("/problemset", isAuthenticated, problemsetView)
	app.Get("/problemset/add", isAuthenticated, addProblemView)
	app.Post("/problemset/add", isAuthenticated, handleAddProblemView)

	app.Get("/problemset/:id", isAuthenticated, showProblemView)
	app.Get("/problemset/:id/dl/:filename", isAuthenticated, downloadProblemInOutFiles)
	app.Get("/problemset/:id/edit", isAuthenticated, editProblemView)
	app.Get("/problemset/:id/:command", isAuthenticated, handlePublishProblemView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
