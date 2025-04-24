package serve

import (
	"fmt"
	"strconv"
	"unicode"
)

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
