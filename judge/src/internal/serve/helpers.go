package serve

import (
	"fmt"
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
