package load_test_data

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/serve"
)

func Erase() {
	connectDatabase()

	var problems []serve.Problem

	if err := db.Where("is_test = ? OR OR owner.is_test = ?", true, true).Find(&problems).Error; err != nil {
		fmt.Println("Error in fetch problems", err)
		return
	}

	for _, problem := range problems {
		folderPath := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))

		if err := os.RemoveAll(folderPath); err != nil {
			fmt.Printf("Error deleting folder %s: %v\n", folderPath, err)
			return
		}

		db.Delete(&problem)
	}
	clear(problems)

	var users []serve.User

	if err := db.Where("is_test = ?", true).Find(&users).Error; err != nil {
		fmt.Println("Error in fetch users", err)
		return
	}

	for _, user := range users {

		db.Delete(&user)
	}
	db.Commit()

}
