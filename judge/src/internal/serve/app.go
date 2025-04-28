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

	go discoverCodeRunners()

	htmlEngine := html.New("static/views/", ".html")
	htmlEngine.AddFunc("add", Add)
	htmlEngine.AddFunc("sub", Sub)
	htmlEngine.AddFunc("mul", Mul)
	htmlEngine.AddFunc("mod", Mod)

	htmlEngine.AddFunc("div", Div)
	htmlEngine.AddFunc("seq", Seq)

	htmlEngine.AddFunc("timeAgo", TimeAgo)

	htmlEngine.AddFunc("truncate", Truncate)
	htmlEngine.AddFunc("breaklines", Breaklines)

	app := fiber.New(fiber.Config{
		Views: htmlEngine,
	})
	app.Static("/static", "static/styles")

	// Landing
	app.Get("/", landingView)

	// login
	app.Get("/login", loginView)
	app.Post("/login", handleLoginView)
	app.Get("/signup", signupView)
	app.Post("/signup", handleSignupView)
	app.Get("/logout", logoutView)

	app.Get("/user/submissions", isAuthenticated, submissionsView)
	app.Get("/user/:id", isAuthenticated, showProfileView)
	app.Get("/user", isAuthenticated, showProfileView)
	app.Post("/user/:id/promote", isAuthenticated, promoteUserView)
	app.Post("/user/:id/demote", isAuthenticated, demoteUserView)

	app.Get("/user/:id/submissions", isAuthenticated, submissionsView)
	app.Get("/user/:id/submissions/dl/:filename", isAuthenticated, downloadSubmissionFiles)

	// problemset
	app.Get("/problemset", isAuthenticated, problemsetView)
	app.Get("/problemset/add", isAuthenticated, addProblemView)
	app.Post("/problemset/add", isAuthenticated, handleAddProblemView)

	app.Get("/problemset/:id", isAuthenticated, showProblemView)
	app.Get("/problemset/:id/dl/:filename", isAuthenticated, downloadProblemInOutFiles)
	app.Get("/problemset/:id/edit", isAuthenticated, editProblemView)
	app.Post("/problemset/:id/edit", isAuthenticated, handleEditProblemView)

	app.Post("/problemset/:id", isAuthenticated, handleSubmitProblemView)

	app.Get("/problemset/:id/:command", isAuthenticated, handlePublishProblemView) // command is publish and unpublish

	app.Post("/code/callback", runCodeCallbackView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}
