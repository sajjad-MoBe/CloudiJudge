package serve

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
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

func Truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}
