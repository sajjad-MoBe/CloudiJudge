package serve

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
)

var store *session.Store

func initSessionStore() {
	store = session.New(
		session.Config{
			Storage: postgres.New(postgres.Config{
				ConnectionURI: fmt.Sprintf(
					"postgres://%s:%s@%s:5432/%s?sslmode=disable",
					os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"),
				),
				Table:      "sessions",
				Reset:      false,
				GCInterval: 10 * time.Minute,
			}),
		},
	)
}

func isAuthenticated(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}
	userID := sess.Get("user_id")
	if userID == nil {
		return c.Redirect("/login")
	}
	c.Locals("user_id", userID)

	return c.Next()
}
