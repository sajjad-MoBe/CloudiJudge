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

func signinView(c *fiber.Ctx) error {
	return render(c, "signin", fiber.Map{
		"Title": "CloudiJudge | ورود کاربر",
	})
}

func landingView(c *fiber.Ctx) error {
	return render(c, "landing", fiber.Map{
		"Title": "CloudiJudge | صفحه اصلی",
	})
}

// Template handling function
func render(c *fiber.Ctx, name string, data interface{}) error {
	return c.Render("pages/"+name, data, "layouts/main")
}
