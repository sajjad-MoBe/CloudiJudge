package serve

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func initSessionStore() {
	store = session.New()
}

func isAuthenticated(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}
	userID := sess.Get("user_id")
	if userID == nil {
		return c.Redirect("/signin")
	}
	c.Locals("user_id", userID)

	return c.Next()
}
