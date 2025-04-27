package serve

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/code_runner"
	"golang.org/x/crypto/bcrypt"
)

// Template handling function
func render(c *fiber.Ctx, name string, data interface{}) error {
	if name[:4] == "sign" {
		return c.Render("pages/"+name, data, "layouts/auth")

	} else if name == "landing" {
		return c.Render("pages/"+name, data, "layouts/main")
	}
	return c.Render("pages/"+name, data, "layouts/dashboard")

}

func error_404(c *fiber.Ctx) error {
	return c.Status(404).Render("pages/error_404", fiber.Map{}, "layouts/main")
}

func error_403(c *fiber.Ctx) error {
	return c.Status(403).Render("pages/error_403", fiber.Map{}, "layouts/main")
}

func loginView(c *fiber.Ctx) error {

	return render(c, "login", fiber.Map{
		"PageTitle": "CloudiJudge | login",
	})
}

func handleLoginView(c *fiber.Ctx) error {

	email := strings.ToLower(c.FormValue("email"))
	password := c.FormValue("password")

	var errorMsg string
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		errorMsg = "Email or password is invalid!"

	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		errorMsg = "Email or password is invalid!"

	} else if sess, err := store.Get(c); err != nil {
		errorMsg = "An unknown error has occurred!"

	} else {
		sess.Set("user_id", user.ID)
		sess.Save()

		return c.Redirect("/problemset")
	}
	return render(c, "login", fiber.Map{
		"PageTitle": "CloudiJudge | login",
		"Message":   errorMsg,
		"Email":     email,
	})

}

func signupView(c *fiber.Ctx) error {
	var email string

	return render(c, "signup", fiber.Map{
		"PageTitle": "CloudiJudge | signup",
		"Email":     email,
	})
}

func handleSignupView(c *fiber.Ctx) error {

	var errorMsg string
	var user User
	// Parse form data
	email := strings.ToLower(c.FormValue("email"))
	password := c.FormValue("password")
	confirm_password := c.FormValue("confirm-password")
	if password != confirm_password {
		errorMsg = "Password and confimation are not same"

	} else if !isSecurePassword(password) {
		errorMsg = "The chosen password must be at least 6 characters long and include letters and numbers."

	} else if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		errorMsg = "Use another password."

	} else if !isValidEmail(email) {
		errorMsg = "The entred email is not valid."

	} else if result := db.Where("email = ?", email).First(&user); result.Error == nil {
		errorMsg = "The entered email is already registered. Please log in."

	} else {

		user = User{
			Email:    email,
			Password: string(hashedPassword),
		}
		db.Create(&user)
		db.Commit()
		return c.Redirect("/login?new=yes")

	}

	return render(c, "signup", fiber.Map{
		"PageTitle": "CloudiJudge | signup",
		"Email":     email,
		"Message":   errorMsg,
	})
}

func logoutView(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err == nil {
		sess.Destroy()
	}

	return c.Redirect("/login")
}

func landingView(c *fiber.Ctx) error {
	return render(c, "landing", fiber.Map{
		"PageTitle": "CloudiJudge | Online Programming Judge",
	})
}

func showProfileView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	userID := uint(id)
	var profileUser User
	result := db.First(&profileUser, userID)
	if result.Error != nil {
		return error_404(c)
	}

	userID = c.Locals("user_id").(uint)
	var user User
	result = db.First(&user, userID)
	if result.Error != nil {
		return error_404(c)
	}
	var submissions []Submission

	err = db.Model(&Submission{}).
		Where("owner_id = ?", profileUser.ID).
		Limit(3).
		Order("created_at DESC").
		Preload("Problem").
		Find(&submissions).Error
	if err != nil {
		clear(submissions)
	}

	return render(c, "show_profile", fiber.Map{
		"PageTitle":   "CloudiJudge | پروفایل کاربر",
		"ProfileUser": profileUser,
		"User":        user,
		"Submissions": submissions,
	})
}

func promoteUserView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	userID := uint(id)
	var targetUser User
	result := db.First(&targetUser, userID)
	if result.Error != nil {
		return error_404(c)
	}

	userID = c.Locals("user_id").(uint)
	var adminUser User
	result = db.First(&adminUser, userID)
	if result.Error == nil {
		if !adminUser.IsAdmin {
			return error_403(c)
		}
		if !targetUser.IsAdmin {
			targetUser.IsAdmin = true
			targetUser.AdminCreatedByID = adminUser.ID
			db.Save(&targetUser)
		}
	}

	return c.Redirect(fmt.Sprintf("/user/%d", targetUser.ID))
}

func demoteUserView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	userID := uint(id)
	var targetUser User
	result := db.First(&targetUser, userID)
	if result.Error != nil {
		return error_404(c)
	}

	userID = c.Locals("user_id").(uint)
	var adminUser User
	result = db.First(&adminUser, userID)
	if result.Error == nil {
		if !adminUser.IsAdmin {
			return error_403(c)
		}
		if targetUser.IsAdmin && targetUser.AdminCreatedByID == adminUser.ID {
			targetUser.IsAdmin = false
			db.Save(&targetUser)
		}
	}

	return c.Redirect(fmt.Sprintf("/user/%d", targetUser.ID))
}

func problemsetView(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(uint)

	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
	}
	limit := c.QueryInt("limit", 10)  // Default limit = 10
	offset := c.QueryInt("offset", 0) // Default offset = 0

	if limit > 100 {
		limit = 100 // Max limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var problems []Problem
	total := publishedProblemsCount

	start := time.Now()

	err := db.Model(&Problem{}).
		Select("id, title, statement, published_at").
		Where("is_published = ?", true).
		Offset(offset).Limit(limit).
		Order("published_at DESC").
		Find(&problems).Error

	duration := time.Since(start)
	fmt.Printf("Time taken: %d ms\n", duration.Milliseconds())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch problems",
		})
	}

	return render(c, "problemset", fiber.Map{
		"Problems": problems,
		"Total":    total,
		"Limit":    limit,
		"Offset":   offset,
		"Pages":    (int(total) + limit - 1) / limit, // Total pages
	})
}

func addProblemView(c *fiber.Ctx) error {

	return render(c, "add_problem", fiber.Map{
		"PageTitle": "CloudiJudge | ساخت سوال",
	})
}

func handleAddProblemView(c *fiber.Ctx) error {
	var errorMsg string = ""
	var problem Problem
	if len(c.FormValue("title")) < 5 || len(c.FormValue("title")) > 50 {
		errorMsg = "طول عنوان وارد شده باید بین 5 الی 50 کاراکتر باشد."

	} else if result := db.Where("title = ?", c.FormValue("title")).First(&problem); result.Error == nil {
		errorMsg = "The selected title is repetitive."

	} else if len(c.FormValue("statement")) < 100 || len(c.FormValue("statement")) > 5000 {
		errorMsg = "طول توضیحات وارد شده باید بین 100 الی 5000 کاراکتر باشد."

	} else if parseInt(c.FormValue("time_limit")) <= 0 {
		errorMsg = "محدودیت زمانی باید یک عدد مثبت باشد."

	} else if parseInt(c.FormValue("memory_limit")) <= 0 {
		errorMsg = "محدودیت حافظه باید یک عدد مثبت باشد."

	} else {

		problem = Problem{
			Title:       c.FormValue("title"),
			Statement:   c.FormValue("statement"),
			TimeLimit:   parseInt(c.FormValue("time_limit")),
			MemoryLimit: parseInt(c.FormValue("memory_limit")),
			OwnerID:     c.Locals("user_id").(uint),
		}

		// Save the problem to the database
		if err := db.Create(&problem).Error; err != nil {
			errorMsg = "خطایی در ذخیره سوال رخ داد."
		} else {
			problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))
			if err := os.MkdirAll(problemDir, os.ModePerm); err != nil {
				errorMsg = "خطایی در ایجاد سوال رخ داد."

			} else if inputFile, err := c.FormFile("input_file"); err != nil {
				errorMsg = "فایل ورودی ها نامعتبر است."

			} else if inputFile.Size > 10*1024*1024 {
				errorMsg = "حجم فایل های ارسالی نباید بیشتر از 10 مگابایت باشد."
			} else if outputFile, err := c.FormFile("output_file"); err != nil {
				errorMsg = "فایل خروجی ها نامعتبر است."

			} else if outputFile.Size > 10*1024*1024 {
				errorMsg = "حجم فایل های ارسالی نباید بیشتر از 10 مگابایت باشد."

			} else if err := c.SaveFile(inputFile, filepath.Join(problemDir, "input.txt")); err != nil {
				errorMsg = "خطایی در ذخیره فایل ورودی ها رخ داد."

			} else if err := c.SaveFile(outputFile, filepath.Join(problemDir, "output.txt")); err != nil {
				errorMsg = "خطایی در ذخیره فایل خروجی ها رخ داد."
			} else {
				return c.Redirect(fmt.Sprintf("/problemset/%d", problem.ID))
			}
			db.Delete(&problem)

		}
	}
	return render(c, "add_problem", fiber.Map{
		"PageTitle":   "CloudiJudge | ساخت سوال",
		"Error":       errorMsg,
		"Title":       c.FormValue("title"),
		"Statement":   c.FormValue("statement"),
		"TimeLimit":   c.FormValue("time_limit"),
		"MemoryLimit": c.FormValue("memory_limit"),
	})
}

func editProblemView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	var problem Problem
	if err := db.Preload("Owner").First(&problem, id).Error; err != nil {
		return error_404(c)
	}

	userID := c.Locals("user_id").(uint)
	if problem.OwnerID != userID {
		var user User
		result := db.First(&user, userID)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
		}
		if !user.IsAdmin {
			return error_403(c)
		}
	}

	return render(c, "add_problem", fiber.Map{
		"PageTitle":   "CloudiJudge | ساخت سوال",
		"Title":       problem.Title,
		"Statement":   problem.Statement,
		"TimeLimit":   problem.TimeLimit,
		"MemoryLimit": problem.MemoryLimit,
		"ProblemID":   problem.ID,
		"Edit":        true,
	})
}

func handleEditProblemView(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	var problem Problem
	if err := db.Preload("Owner").First(&problem, id).Error; err != nil {
		return error_404(c)
	}

	userID := c.Locals("user_id").(uint)
	if problem.OwnerID != userID {
		var user User
		result := db.First(&user, userID)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
		}
		if !user.IsAdmin {
			return error_403(c)
		}
	}

	var errorMsg string = ""

	if len(c.FormValue("title")) < 5 || len(c.FormValue("title")) > 50 {
		errorMsg = "طول عنوان وارد شده باید بین 5 الی 50 کاراکتر باشد."

	} else if len(c.FormValue("statement")) < 100 || len(c.FormValue("statement")) > 5000 {
		errorMsg = "طول توضیحات وارد شده باید بین 100 الی 5000 کاراکتر باشد."

	} else if parseInt(c.FormValue("time_limit")) <= 0 {
		errorMsg = "محدودیت زمانی باید یک عدد مثبت باشد."

	} else if parseInt(c.FormValue("memory_limit")) <= 0 {
		errorMsg = "محدودیت حافظه باید یک عدد مثبت کمتر از هزار باشد."

	} else if parseInt(c.FormValue("memory_limit")) >= 1000 {
		errorMsg = "محدودیت حافظه باید یک عدد مثبت کمتر از هزار باشد."

	} else {

		problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))
		inputFile, inputErr := c.FormFile("input_file")
		outputFile, outputErr := c.FormFile("output_file")

		if inputErr == nil {
			if inputFile.Size > 10*1024*1024 {
				errorMsg = "حجم فایل های ارسالی نباید بیشتر از 10 مگابایت باشد."
			} else if err := c.SaveFile(inputFile, filepath.Join(problemDir, "input.txt")); err != nil {
				errorMsg = "خطایی در ذخیره فایل ورودی ها رخ داد."
			}
		}
		if errorMsg == "" && outputErr == nil {
			if outputFile.Size > 10*1024*1024 {
				errorMsg = "حجم فایل های ارسالی نباید بیشتر از 10 مگابایت باشد."
			} else if err := c.SaveFile(outputFile, filepath.Join(problemDir, "output.txt")); err != nil {
				errorMsg = "خطایی در ذخیره فایل خروجی ها رخ داد."
			}
		}
		if errorMsg == "" {
			problem.Title = c.FormValue("title")
			problem.Statement = c.FormValue("statement")
			problem.TimeLimit = parseInt(c.FormValue("time_limit"))
			problem.MemoryLimit = parseInt(c.FormValue("memory_limit"))
			problem.IsPublished = false
			db.Save(&problem)
			return c.Redirect(fmt.Sprintf("/problemset/%d", problem.ID))
		}

	}
	return render(c, "add_problem", fiber.Map{
		"PageTitle":   "CloudiJudge | ساخت سوال",
		"Error":       errorMsg,
		"Title":       c.FormValue("title"),
		"Statement":   c.FormValue("statement"),
		"TimeLimit":   c.FormValue("time_limit"),
		"MemoryLimit": c.FormValue("memory_limit"),
		"ProblemID":   problem.ID,
		"Edit":        true,
	})
}

func showProblemView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	var problem Problem
	if err := db.Preload("Owner").First(&problem, id).Error; err != nil {
		return error_404(c)
	}

	userID := c.Locals("user_id")
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
	}

	// db.Save(&user)
	return render(c, "show_problem", fiber.Map{
		"PageTitle": "CloudiJudge | مشاهده سوال",
		"Problem":   problem,
		"User":      user,
	})
}

func downloadProblemInOutFiles(c *fiber.Ctx) error {

	filename := c.Params("filename")
	if filename != "input.txt" && filename != "output.txt" {
		return error_404(c)
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	var problem Problem
	if err := db.Preload("Owner").First(&problem, id).Error; err != nil {
		return error_404(c)
	}

	userID := c.Locals("user_id")
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
	}
	if !problem.IsPublished && problem.OwnerID != user.ID && !user.IsAdmin {
		return error_403(c)
	}

	filePath := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d/%s", problem.ID, filename))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return error_404(c)
	}

	return c.Download(filePath, filename)
}

func handlePublishProblemView(c *fiber.Ctx) error {

	command := c.Params("command")
	if command != "publish" && command != "unpublish" {
		return error_404(c)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	var problem Problem
	if err := db.First(&problem, id).Error; err != nil {
		return error_404(c)
	}

	userID := c.Locals("user_id")
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
	}
	if !user.IsAdmin {
		return error_403(c) // change to no permission
	}

	if command == "publish" {
		if !problem.IsPublished {
			problem.IsPublished = true
			now := time.Now()
			problem.PublishedAt = &now
		}
		publishedProblemsCount++
	} else {
		publishedProblemsCount--
		problem.IsPublished = false
	}
	db.Save(&problem)

	return c.Redirect(fmt.Sprintf("/problemset/%d", problem.ID))
}

func handleSubmitProblemView(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	var problem Problem
	if err := db.Preload("Owner").First(&problem, id).Error; err != nil {
		return error_404(c)
	}
	if !problem.IsPublished {
		return error_403(c)
	}

	userID := c.Locals("user_id").(uint)
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
	}
	submission := Submission{
		Status:    "waiting",
		Token:     GenerateRandomToken(30),
		OwnerID:   user.ID,
		ProblemID: problem.ID,
	}
	var errorMsg string
	// Save the problem to the database
	if err := db.Create(&submission).Error; err != nil {
		errorMsg = "خطایی در ذخیره ارسال رخ داد."

	} else {
		problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))

		if submittedFile, err := c.FormFile("submit_file"); err != nil {
			fmt.Println(err)
			errorMsg = "فایل ارسالی نامعتبر است."

		} else if ext := filepath.Ext(submittedFile.Filename); ext != ".go" {
			errorMsg = "فقط فایل های گولنگ قابل قبول هستند."

		} else if submittedFile.Size > 10*1024*1024 {
			errorMsg = "حجم فایل ارسالی نباید بیشتر از 10 مگابایت باشد."

		} else if err := c.SaveFile(submittedFile,
			filepath.Join(problemDir, strconv.Itoa(int(submission.ID))+".go")); err != nil {

			errorMsg = "خطایی در ذخیره فایل ارسالی رخ داد."

		} else {
			user.SolveAttemps++
			db.Save(&user)
			go sendCodeToRun(submission, problem)
			return c.Redirect(fmt.Sprintf("/user/%d/submissions", user.ID))
		}
	}
	db.Delete(&submission)
	return render(c, "show_problem", fiber.Map{
		"PageTitle": "CloudiJudge | مشاهده سوال",
		"Error":     errorMsg,
		"Problem":   problem,
		"User":      user,
	})
}

func submissionsView(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}
	userID := uint(id)
	var targetUser User
	result := db.First(&targetUser, userID)
	if result.Error != nil {
		fmt.Println(result.Error)
		return error_404(c)
	}

	userID = c.Locals("user_id").(uint)
	var thisUser User
	result = db.First(&thisUser, userID)
	if thisUser.ID != targetUser.ID {
		if result.Error == nil {
			if !thisUser.IsAdmin {
				return error_403(c)
			}
		}
	}

	limit := c.QueryInt("limit", 10)  // Default limit = 10
	offset := c.QueryInt("offset", 0) // Default offset = 0

	if limit > 100 {
		limit = 100 // Max limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var submissions []Submission
	var total int64

	err = db.Model(&Submission{}).
		Where("owner_id = ?", targetUser.ID).
		Count(&total).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Preload("Problem").
		Find(&submissions).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch submissions",
		})
	}

	return render(c, "show_submissions", fiber.Map{
		"Submissions": submissions,
		"ProfileUser": targetUser,
		"Total":       int(total),
		"Limit":       limit,
		"Offset":      offset,
		"Pages":       (int(total) + limit - 1) / limit, // Total pages
	})
}

func downloadSubmissionFiles(c *fiber.Ctx) error {

	filename, err := strconv.Atoi(c.Params("filename"))
	if err != nil {
		return error_404(c)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return error_404(c)
	}

	var submission Submission
	if err := db.First(&submission, filename).Error; err != nil || id != int(submission.OwnerID) {
		return error_404(c)
	}

	userID := c.Locals("user_id").(uint)
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("خطا در دریافت اطلاعات کاربر")
	}
	if submission.OwnerID != userID && !user.IsAdmin {
		return error_403(c)
	}

	filePath := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d/%d.go", submission.ProblemID, filename))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return error_404(c)
	}

	return c.Download(filePath, "main.go")
}

func runCodeCallbackView(c *fiber.Ctx) error {
	var data code_runner.ResultData

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("")
	}
	var submission Submission
	if err := db.Where("token = ?", data.CallbackToken).Preload("Owner").First(&submission).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("")
	}
	submission.Status = data.Status
	if data.Status == "Accepted" {
		user := submission.Owner
		user.SuccessAttemps += 1
		db.Save(&user)
	}
	db.Save(&submission)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok": true,
	})
}
