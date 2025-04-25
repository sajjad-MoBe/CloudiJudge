package serve

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
	"unicode"
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
func Div(v1, v2 int) int {
	if v2 == 0 {
		return 0
	}
	return v1 / v2
}

func Truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}

func TimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Hour {
		minutes := int(duration.Minutes())
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
