package serve

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

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
	var message string
	sess, err := store.Get(c)
	if err == nil {
		message = sess.Get("siginError").(string)
		sess.Delete("siginError")
		sess.Save()
	}

	return render(c, "signin", fiber.Map{
		"Title":   "CloudiJudge | ورود کاربر",
		"Message": message,
	})
}

func signinSubmitView(c *fiber.Ctx) error {

	// Parse form data
	email := c.FormValue("email")
	password := c.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		setSigninError(c, "ایمیل یا رمز عبور وارد شده معتبر نمیباشد.")
		c.Redirect("/signin")
	}
	password = string(hashedPassword)

	// Validate user credentials
	var user User
	result := db.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		setSigninError(c, "ایمیل یا رمز عبور وارد شده معتبر نمیباشد.")
		return c.Redirect("/signin")
	}

	sess, err := store.Get(c)
	if err != nil {
		setSigninError(c, "خطای ناشناخته ای رخ داد")
		c.Redirect("/signin")
	}
	sess.Set("user_id", user.ID)
	sess.Save()

	// Redirect to the dashboard
	return c.Redirect("/dashboard")
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

func setSigninError(c *fiber.Ctx, errMsg string) {
	sess, err := store.Get(c)
	if err == nil {
		sess.Set("siginError", errMsg)
		sess.Save()
	}
}
