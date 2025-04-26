package serve

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/sajjad-MoBe/CloudiJudge/judge/src/internal/code_runner"
)

func GenerateRandomToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "sfjskfjejowiri3938*@(#&E489643)"
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func isSecurePassword(password string) bool {
	if len(password) <= 6 {
		return false
	}

	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsDigit(char) {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}

func parseFloat32(value string) float32 {
	if value == "" {
		return 0
	}
	var result float32
	fmt.Sscanf(value, "%f", &result)
	return result
}

func parseInt(value string) int {
	ress, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return ress
}

func Add(a, b int) int {
	return a + b
}

func Sub(a, b int) int {
	return a - b
}
func Mul(a, b int) int {
	return a * b
}
func Div(v1, v2 int) int {
	if v2 == 0 {
		return 0
	}
	return v1 / v2
}
func Mod(v1, v2 int) int {
	if v2 == 0 {
		return 0
	}
	return v1 % v2
}

func Seq(start, end int) []int {
	result := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		result[i-start] = i
	}
	return result
}

func Truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}

func TimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < 2*time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 0 {
			return "just now"
		}
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)

	} else if duration < 48*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)

	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

func sendCodeToRun(submission Submission, problem Problem) {
	run := code_runner.Run{
		TimeLimitMs:   problem.TimeLimit,
		PproblemID:    int(problem.ID),
		SubmissionID:  int(submission.ID),
		CallbackToken: submission.Token,
	}
	jsonData, err := json.Marshal(run)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// dynamic code runner
	resp, err := http.Post("http://localhost:2/run", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		submission.Status = "Compilation failed"
		db.Save(&submission)
		return
	}

}
