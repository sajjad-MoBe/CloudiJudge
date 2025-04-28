package serve

import (
	"fmt"
	"log"
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

func error_404(c *fiber.Ctx, messages ...string) error {
	var message string
	if len(messages) > 0 {
		message = messages[0]
	}
	return c.Status(404).Render("pages/error_404", fiber.Map{
		"PageTitle": "Page Not Found",
		"Message":   message,
	}, "layouts/error")
}

func error_403(c *fiber.Ctx, messages ...string) error {
	var message string
	if len(messages) > 0 {
		message = messages[0]
	}
	return c.Status(403).Render("pages/error_403", fiber.Map{
		"PageTitle": "Access Denied",
		"Message":   message,
	}, "layouts/error")
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
		errorMsg = "An unknown error has been occurred!"

	} else {
		sess.Set("user_id", user.ID)
		sess.Save()
		return c.Redirect("/problemset")
	}
	return render(c.Status(fiber.StatusBadRequest), "login", fiber.Map{
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
		db.Save(&user)

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
	userID := c.Locals("user_id").(uint)
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return c.Redirect("/login")
	}

	var profileUser User
	if c.Path() == "/user" {
		profileUser = user
	} else {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return error_404(c)
		}
		userID = uint(id)
		result = db.First(&profileUser, userID)
		if result.Error != nil {
			return error_404(c)
		}
	}

	var submissions []Submission

	err := db.Model(&Submission{}).
		Where("owner_id = ?", profileUser.ID).
		Limit(3).
		Order("created_at DESC").
		Preload("Problem").
		Find(&submissions).Error
	if err != nil {
		clear(submissions)
	}

	return render(c, "show_profile", fiber.Map{
		"PageTitle":   "CloudiJudge | profile",
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
		return c.Redirect("/login")
	}
	myproblems := c.Query("myproblems", "-")
	limit := c.QueryInt("limit", 10)  // Default limit = 10
	offset := c.QueryInt("offset", 0) // Default offset = 0

	if limit > 100 {
		limit = 100 // Max limit = 100
	} else if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}

	var problems []Problem
	total := int64(publishedProblemsCount)
	var err error
	// start := time.Now()
	if myproblems == "yes" {
		if user.IsAdmin {
			total = int64(publishedProblemsCount + notPublishedProblemsCount)
			err = db.Model(&Problem{}).
				Select("id, title, statement, is_published, published_at").
				// Where("is_published = ?", true).
				Offset(offset).Limit(limit).
				Order("published_at DESC").
				Find(&problems).Error

		} else {
			err = db.Model(&Problem{}).
				Select("id, title, statement, is_published, published_at").
				Where("owner_id = ?", user.ID).
				Count(&total).
				Offset(offset).Limit(limit).
				Order("published_at DESC").
				Find(&problems).Error
		}
	} else {
		myproblems = "no"
		err = db.Model(&Problem{}).
			Select("id, title, statement, is_published, published_at").
			Where("is_published = ?", true).
			Offset(offset).Limit(limit).
			Order("published_at DESC").
			Find(&problems).Error
	}

	// duration := time.Since(start)
	// fmt.Printf("Time taken: %d ms\n", duration.Milliseconds())

	if err != nil {
		return render(c.Status(fiber.StatusInternalServerError), "problemset", fiber.Map{
			"PageTitle":  "CloudiJudge | problemset",
			"User":       user,
			"Problems":   problems,
			"Total":      int(total),
			"Limit":      0,
			"Offset":     0,
			"Pages":      0,
			"Myproblems": myproblems,
		})
	}

	return render(c, "problemset", fiber.Map{
		"PageTitle":   "CloudiJudge | problemset",
		"User":        user,
		"Problems":    problems,
		"Total":       int(total),
		"Limit":       limit,
		"Offset":      offset,
		"CurrentPage": (offset / limit) + 1,
		"Pages":       (int(total) + limit - 1) / limit, // Total pages
		"Myproblems":  myproblems,
	})
}

func addProblemView(c *fiber.Ctx) error {

	return render(c, "add_problem", fiber.Map{
		"PageTitle": "CloudiJudge | add problem",
	})
}

func handleAddProblemView(c *fiber.Ctx) error {
	var errorMsg string = ""
	var problem Problem
	if len(c.FormValue("title")) < 5 || len(c.FormValue("title")) > 50 {
		errorMsg = "Title length most be between 5 and 50."

	} else if result := db.Where("title = ?", c.FormValue("title")).First(&problem); result.Error == nil {
		errorMsg = "The selected title is repetitive."

	} else if len(c.FormValue("statement")) < 100 || len(c.FormValue("statement")) > 5000 {
		errorMsg = "Statement length most be between 100 and 5000."

	} else if parseInt(c.FormValue("time_limit")) <= 0 {
		errorMsg = "Time limit most be a positive number."

	} else if parseInt(c.FormValue("memory_limit")) <= 0 {
		errorMsg = "Memory limit most be a positive number."

	} else if parseInt(c.FormValue("memory_limit")) <= 500 {
		errorMsg = "Memory limit most be less than 500."

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
			log.Println("error in create problem", err)
			errorMsg = "An unknown error has been occurred!"
		} else {
			problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))
			if err := os.MkdirAll(problemDir, os.ModePerm); err != nil {
				log.Println("error in make directory for new problem", err)
				errorMsg = "An unknown error has been occurred!"

			} else if inputFile, err := c.FormFile("input_file"); err != nil {
				errorMsg = "Invalid input file."

			} else if inputFile.Size > 10*1024*1024 {
				errorMsg = "Input file size most be less than 10Mb"

			} else if outputFile, err := c.FormFile("output_file"); err != nil {
				errorMsg = "Invalid output file."

			} else if outputFile.Size > 10*1024*1024 {
				errorMsg = "Output file size most be less than 10Mb"

			} else if err := c.SaveFile(inputFile, filepath.Join(problemDir, "input.txt")); err != nil {
				log.Println("error in save input file for new problem", err)
				errorMsg = "An unknown error has been occurred!"

			} else if err := c.SaveFile(outputFile, filepath.Join(problemDir, "output.txt")); err != nil {
				log.Println("error in save output file for new problem", err)
				errorMsg = "An unknown error has been occurred!"
			} else {
				return c.Redirect(fmt.Sprintf("/problemset/%d", problem.ID))
			}
			db.Delete(&problem)

		}
	}
	return render(c, "add_problem", fiber.Map{
		"PageTitle":   "CloudiJudge | add problem",
		"Message":     errorMsg,
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
			return c.Redirect("/login")
		}
		if !user.IsAdmin {
			return error_403(c)
		}
	}

	return render(c, "add_problem", fiber.Map{
		"PageTitle":   "CloudiJudge | edit problem",
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
			return c.Redirect("/login")
		}
		if !user.IsAdmin {
			return error_403(c)
		}
	}

	var errorMsg string = ""

	if len(c.FormValue("title")) < 5 || len(c.FormValue("title")) > 50 {
		errorMsg = "Title length most be between 5 and 50."

	} else if len(c.FormValue("statement")) < 100 || len(c.FormValue("statement")) > 5000 {
		errorMsg = "Statement length most be between 100 and 5000."

	} else if parseInt(c.FormValue("time_limit")) <= 0 {
		errorMsg = "Time limit most be a positive number."

	} else if parseInt(c.FormValue("memory_limit")) <= 0 {
		errorMsg = "Memory limit most be a positive number."

	} else if parseInt(c.FormValue("memory_limit")) >= 500 {
		errorMsg = "Memory limit most be less than 500."

	} else {
		problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))
		inputFile, inputErr := c.FormFile("input_file")
		outputFile, outputErr := c.FormFile("output_file")

		if inputErr == nil {
			if inputFile.Size > 10*1024*1024 {
				errorMsg = "Input file size most be less than 10Mb"

			} else if err := c.SaveFile(inputFile, filepath.Join(problemDir, "input.txt")); err != nil {
				log.Println("error in save input file for problem", problem.ID, err)
				errorMsg = "An unknown error has been occurred!"
			}
		}
		if errorMsg == "" && outputErr == nil {
			if outputFile.Size > 10*1024*1024 {
				errorMsg = "Output file size most be less than 10Mb"

			} else if err := c.SaveFile(outputFile, filepath.Join(problemDir, "output.txt")); err != nil {
				log.Println("error in save output file for problem", problem.ID, err)
				errorMsg = "An unknown error has been occurred!"
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
		"PageTitle":   "CloudiJudge | edit problem",
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
		return c.Redirect("/login")
	}

	// db.Save(&user)
	return render(c, "show_problem", fiber.Map{
		"PageTitle": "CloudiJudge | view problem",
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
		return c.Redirect("/login")
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
		return c.Redirect("/login")
	}
	if !user.IsAdmin {
		return error_403(c)
	}

	if command == "publish" {
		if !problem.IsPublished {
			problem.IsPublished = true
			now := time.Now()
			problem.PublishedAt = &now
		}
		publishedProblemsCount++
		notPublishedProblemsCount--
	} else {
		publishedProblemsCount--
		notPublishedProblemsCount++
		problem.IsPublished = false
	}
	db.Save(&problem)
	if c.Query("next", "-") == "problemset" {
		if c.Query("myproblems", "-") == "yes" {
			return c.Redirect("/problemset?myproblems=yes")
		}
		return c.Redirect("/problemset")
	}
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
		return c.Redirect("/login")
	}
	submission := Submission{
		Status:    "waiting",
		Token:     GenerateRandomToken(30),
		OwnerID:   user.ID,
		ProblemID: problem.ID,
	}
	var errorMsg string

	if err := db.Create(&submission).Error; err != nil {
		log.Println("error in save new submission", err)
		errorMsg = "An unknown error has been occurred!"

	} else {
		problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))

		if submittedFile, err := c.FormFile("submit_file"); err != nil {
			errorMsg = "Invalid submission code file!"

		} else if ext := filepath.Ext(submittedFile.Filename); ext != ".go" {
			errorMsg = "Submission code is not GO!"

		} else if submittedFile.Size > 10*1024*1024 {
			errorMsg = "Submission code size most be lower than 10Mb!"

		} else if err := c.SaveFile(submittedFile,
			filepath.Join(problemDir, strconv.Itoa(int(submission.ID))+".go")); err != nil {

			log.Println("error in save submission code file", err)
			errorMsg = "An unknown error has been occurred!"

		} else {
			user.SolveAttemps++
			db.Save(&user)
			go sendCodeToRun(submission, problem)
			return c.Redirect(fmt.Sprintf("/user/%d/submissions", user.ID))
		}
	}
	db.Delete(&submission)
	return render(c, "show_problem", fiber.Map{
		"PageTitle": "CloudiJudge | show problem",
		"Message":   errorMsg,
		"Problem":   problem,
		"User":      user,
	})
}

func submissionsView(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var thisUser User
	result := db.First(&thisUser, userID)
	if result.Error != nil {
		return c.Redirect("/login")
	}

	var targetUser User
	if c.Path() == "/user/submissions" {
		targetUser = thisUser
	} else {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return error_404(c)
		}
		userID = uint(id)

		result = db.First(&targetUser, userID)
		if result.Error != nil {
			fmt.Println(result.Error)
			return error_404(c)
		}
	}

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

	err := db.Model(&Submission{}).
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
		"CurrentPage": (offset / limit) + 1,
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
		return error_404(c)
	}
	var submission Submission
	if err := db.Where("token = ?", data.CallbackToken).Preload("Owner").First(&submission).Error; err != nil {
		return error_404(c)
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
