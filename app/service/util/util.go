package util

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string, hashedPassword string) bool {
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func IsValidEmail(email string) bool {
	// Define a regular expression for email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Match the email against the regular expression
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func IsValidPhone(phone string) bool {
	// Remove any non-digit characters from the phone number string
	phone = strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// Check the length of the cleaned phone number
	if len(phone) >= 10 && len(phone) <= 13 {
		// If the phone number starts with "+62," replace it with "08"
		if strings.HasPrefix(phone, "+62") {
			phone = "08" + phone[3:]
		}

		// Ensure the phone number starts with "08"
		if strings.HasPrefix(phone, "08") {
			// Check if the resulting phone number contains only digits
			if _, err := strconv.Atoi(phone); err == nil {
				return true
			}
		}
	}
	return false
}

func GenerateHash(password string) (string, error) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func Int(v int) *int { return &v }
