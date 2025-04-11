package serve

import "github.com/gofiber/fiber/v2"

func problemsetView(c *fiber.Ctx) error {

	// Retrieve user ID from the session
	userID := c.Locals("user_id")

	// Fetch user details from the database
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
	}
	return render(c, "problemset", fiber.Map{
		"Title": "CloudiJudge | سوالات",
	})
}
