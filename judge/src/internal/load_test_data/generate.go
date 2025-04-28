package load_test_data

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/serve"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAndFill() {
	connectDatabase()
	testPassword, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Data can't be generated")
		return
	}
	for i := 1; i < 10_001; i++ {
		fmt.Println("adding user", i+1)
		user := serve.User{
			Email:    fmt.Sprintf("test_user_%d@gmail.com", i),
			Password: string(testPassword),
			IsTest:   true,
		}
		db.Create(&user)
	}
	fmt.Println("test users were generated")
	for i := 0; i < 50_000; i++ {
		fmt.Println("adding problem", i+1)
		randomIndex := rand.Intn(50_000) + 1
		isPublished := randomIndex%2 == 0
		randomTime := time.Now().Add(-time.Duration(rand.Int63n(5*60*60)) * time.Second)

		problem := serve.Problem{
			Title:       fmt.Sprintf("test problem %d", i+1),
			Statement:   generateLoremIpsum(rand.Int()%30 + 20),
			IsPublished: isPublished,
			PublishedAt: &randomTime,
			TimeLimit:   rand.Int()%10 + 1,
			MemoryLimit: rand.Int()%10 + 1,
			OwnerID:     uint(randomIndex),
			IsTest:      true,
		}
		db.Create(&problem)

		problemDir := filepath.Join(os.Getenv("PROBLEM_UPLOAD_FOLDER"), fmt.Sprintf("%d", problem.ID))
		if err := os.MkdirAll(problemDir, os.ModePerm); err != nil {
			db.Delete(&problem)
			fmt.Println("failed to generate test problem")
			return
		}
		file, err := os.Create(filepath.Join(problemDir, "input.txt"))
		if err != nil {
			fmt.Println("Error creating input file:", err)
			return
		}
		_, err = file.WriteString("1\n2\n3")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			file.Close()
			return
		}
		file.Close()
		file, err = os.Create(filepath.Join(problemDir, "output.txt"))
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		_, err = file.WriteString("2\n4\n6")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			file.Close()
			return
		}
		file.Close()
	}
	db.Commit()

}

func generateLoremIpsum(wordCount int) string {
	loremWords := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
		"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
		"magna", "aliqua", "ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
		"exercitation", "ullamco", "folan", "nisi", "ut", "aliquip", "chetorii", "ea",
		"commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
		"in", "voluptate", "velit", "salam", "cillum", "dolore", "eu", "fugiat", "nulla",
		"pariatur", "excepteur", "sint", "occaecat", "ali", "non", "proident",
		"sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id",
		"est", "laborum",
	}

	var sb strings.Builder
	for i := 0; i < wordCount; i++ {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(loremWords[rand.Intn(len(loremWords))])
	}
	return sb.String()
}
