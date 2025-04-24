package serve

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Template handling function
func render(c *fiber.Ctx, name string, data interface{}) error {
	return c.Render("pages/"+name, data, "layouts/main")
}

func setSigninError(c *fiber.Ctx, email, errMsg string) {
	sess, err := store.Get(c)
	if err == nil {
		sess.Set("siginError", errMsg)
		sess.Set("email", email)
		sess.Save()
	}
}

func signinView(c *fiber.Ctx) error {
	var email string
	var message string

	sess, err := store.Get(c)
	if err == nil {
		tmp := sess.Get("siginError")
		if tmp != nil {
			message = tmp.(string)
			sess.Delete("siginError")
			sess.Save()
		}
		tmp = sess.Get("email")
		if tmp != nil {
			email = tmp.(string)
			sess.Delete("email")
			sess.Save()
		}
	}

	return render(c, "signin", fiber.Map{
		"PageTitle": "CloudiJudge | ورود کاربر",
		"Message":   message,
		"Email":     email,
	})
}

func handleSigninView(c *fiber.Ctx) error {

	// Parse form data
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Validate user credentials
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		setSigninError(c, email, "ایمیل یا رمز عبور وارد شده معتبر نمیباشد.")
		return c.Redirect("/signin")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		setSigninError(c, email, "ایمیل یا رمز عبور وارد شده معتبر نمیباشد.")
	}

	sess, err := store.Get(c)
	if err != nil {
		setSigninError(c, email, "خطای ناشناخته ای رخ داد")
		c.Redirect("/signin")
	}
	sess.Set("user_id", user.ID)
	sess.Save()

	// Redirect to the problemset
	return c.Redirect("/problemset")
}

func handleSignupView(c *fiber.Ctx) error {

	// Parse form data
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirm_password := c.FormValue("confirm_password")

	if password != confirm_password {
		setSigninError(c, email, "رمز عبور و تایید آن یکسان نیستند.")
		c.Redirect("/signin")
	}

	if !isSecurePassword(password) {
		setSigninError(c, email, "رمز عبور انتخابی باید حداقل به طول ۶ و شامل حروف و اعداد باشد.")
		c.Redirect("/signin")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		setSigninError(c, email, "رمز عبور دیگری انتخاب کنید.")
		c.Redirect("/signin")
	}
	password = string(hashedPassword)

	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error == nil {
		setSigninError(c, email, "ایمیل وارد شده قبلا ثبت نام کرده است. لطفا وارد شوید.")
		return c.Redirect("/signin")
	}
	db.Create(&User{
		Email:    email,
		Password: password,
	})
	db.Commit()
	sess, err := store.Get(c)
	if err != nil {
		setSigninError(c, email, "خطای ناشناخته ای رخ داد")
		c.Redirect("/signin")
	}
	sess.Set("user_id", user.ID)
	sess.Save()

	// Redirect to the problemset
	return c.Redirect("/problemset")
}

func signoutView(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err == nil {
		sess.Destroy()
	}

	return c.Redirect("/signin")
}

func landingView(c *fiber.Ctx) error {
	return render(c, "landing", fiber.Map{
		"PageTitle": "CloudiJudge | صفحه اصلی",
	})
}

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
		"PageTitle": "CloudiJudge | سوالات",
	})
}

func addProblemView(c *fiber.Ctx) error {

	return render(c, "add_problem", fiber.Map{
		"PageTitle": "CloudiJudge | ساخت سوال",
	})
}

func handleAddProblemView(c *fiber.Ctx) error {
	var errorMsg string = ""

	if c.FormValue("Pagetitle") == "" {
		errorMsg = "عنوان وارد شده نامعتبر است."

	} else if c.FormValue("content") == "" {
		errorMsg = "توضیحات وارد شده نامعتبر است."

	} else if parseFloat32(c.FormValue("time_limit")) <= 0 {
		errorMsg = "محدودیت زمانی باید یک عدد مثبت باشد."

	} else if parseFloat32(c.FormValue("memory_limit")) <= 0 {
		errorMsg = "محدودیت حافظه باید یک عدد مثبت باشد."

	} else {

		problem := Problem{
			Title:       c.FormValue("Pagetitle"),
			Content:     c.FormValue("content"),
			TimeLimit:   parseFloat32(c.FormValue("time_limit")),
			MemoryLimit: parseFloat32(c.FormValue("memory_limit")),
			UserID:      c.Locals("user_id").(uint),
		}

		// Save the problem to the database
		if err := db.Create(&problem).Error; err != nil {
			fmt.Println(err)
			errorMsg = "خطایی در ذخیره سوال رخ داد."
		} else {
			problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))
			if err := os.MkdirAll(problemDir, os.ModePerm); err != nil {
				errorMsg = "خطایی در ایجاد سوال رخ داد."

			} else if inputFile, err := c.FormFile("input_file"); err != nil {
				errorMsg = "فایل ورودی ها نامعتبر است."

			} else if outputFile, err := c.FormFile("output_file"); err != nil {
				errorMsg = "فایل خروجی ها نامعتبر است."

			} else if err := c.SaveFile(inputFile, filepath.Join(problemDir, "input.txt")); err != nil {
				fmt.Println(err)
				errorMsg = "خطایی در ذخیره فایل ورودی ها رخ داد."

			} else if err := c.SaveFile(outputFile, filepath.Join(problemDir, "output.txt")); err != nil {
				fmt.Println(err)
				errorMsg = "خطایی در ذخیره فایل خروجی ها رخ داد."
			} else {
				return c.Redirect(fmt.Sprintf("/problemset/%d", problem.ID))
			}

		}
	}

	return render(c, "add_problem", fiber.Map{
		"PageTitle":   "CloudiJudge | ساخت سوال",
		"Error":       errorMsg,
		"Title":       c.FormValue("title"),
		"Content":     c.FormValue("content"),
		"TimeLimit":   c.FormValue("time_limit"),
		"MemoryLimit": c.FormValue("memory_limit"),
	})
}

func showProblemView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid problem ID"})
	}

	var problem Problem
	if err := db.First(&problem, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Problem not found"})
	}

	return render(c, "show_problem", fiber.Map{
		"PageTitle":   "CloudiJudge | مشاهده سوال",
		"Title":       problem.Title,
		"Content":     problem.Content,
		"TimeLimit":   problem.TimeLimit,
		"MemoryLimit": problem.MemoryLimit,
		"InputFile":   "problem.InputFile",
		"OutputFile":  "problem.OutputFile",
	})
}
