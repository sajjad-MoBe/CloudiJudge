package serve

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// Template handling function
func render(c *fiber.Ctx, name string, data interface{}) error {
	return c.Render("pages/"+name, data, "layouts/main")
}

func StartListening(port int) {
	connectDatabase()
	initSessionStore()

	app := fiber.New(fiber.Config{
		Views: html.New("static/views/", ".html"),
	})

	// Landing
	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, "landing", fiber.Map{
			"Title": "CloudiJudge | صفحه اصلی",
		})
	})

	// signin
	app.Get("/signin", func(c *fiber.Ctx) error {
		return render(c, "signin", fiber.Map{
			"Title": "CloudiJudge | ورود کاربر",
		})
	})

	// problemset
	app.Get("/problemset", isAuthenticated, problemsetView)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
