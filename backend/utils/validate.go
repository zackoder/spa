package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsValidName(w http.ResponseWriter, name string) bool {
	if len(name) == 0 {
		// CreateJson(w, "Enter your information", http.StatusBadRequest)
		return false
	}
	if len(name) > 20 {
		// CreateJson(w, "should have 8 to 20 Letters", http.StatusBadRequest)
		return false
	}

	reg := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{2,20}$`)
	if !reg.MatchString(name) {
		// CreateJson(w, "enter your information correctly", http.StatusBadRequest)
		return false
	}
	return true
}

func IsValidGender(w http.ResponseWriter, gender string) bool {
	if gender == "" || (gender != "male" && gender != "female") {
		// CreateJson(w, "enter your gender", http.StatusBadRequest)
		return false
	}

	return true
}

func IsValidAge(w http.ResponseWriter, age string) bool {
	// let dateRegex = /^\d{4}-\d{2}-\d{2}$/;
	reg := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !reg.MatchString(age) {
		CreateJson(w, "Please Check your Birth day", http.StatusBadRequest)
		return false
	}

	// Split the date string into year, month, and day
	parts := strings.Split(age, "-")
	if len(parts) != 3 {
		CreateJson(w, "Invalid date format", http.StatusBadRequest)
		return false
	}

	// Convert the parts to integers
	year, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	day, err3 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil || err3 != nil {
		CreateJson(w, "Invalid date format", http.StatusBadRequest)
		return false
	}

	// Create a time.Time object for the birth date
	birthDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Validate the date
	isValidDate := birthDate.Year() == year && birthDate.Month() == time.Month(month) && birthDate.Day() == day

	if !isValidDate {
		CreateJson(w, "Invalid date format", http.StatusBadRequest)
		return false
	}

	// Get the current date
	currentDate := time.Now()

	// Calculate the age difference
	ageDiff := currentDate.Year() - birthDate.Year()

	// Check if the birthday has occurred this year
	hasBirthdayOccurred := currentDate.Month() > birthDate.Month() ||
		(currentDate.Month() == birthDate.Month() && currentDate.Day() >= birthDate.Day())

	if !hasBirthdayOccurred {
		ageDiff--
	}

	// Check if the age is at least 10
	if ageDiff >= 10 {
		return true
	} else {
		CreateJson(w, "your age not accessible", http.StatusBadRequest)
		return false
	}
}

func IsValidEmail(w http.ResponseWriter, email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if email == "" {
		CreateJson(w, "Please enter your email", http.StatusBadRequest)
		return false
	}
	if !emailRegex.MatchString(email) {

		CreateJson(w, "Please enter your email correctly", http.StatusBadRequest)
		return false
	}
	return true
}

func IsValidPassword(w http.ResponseWriter, password string) bool {
	if password == "" {
		CreateJson(w, "Please enter your password", http.StatusBadRequest)
		return false
	}

	if len(password) < 8 {
		CreateJson(w, "Please enter at least 8 letters", http.StatusBadRequest)
		return false
	}

	reg := regexp.MustCompile(`[a-z]`)
	if !reg.MatchString(password) {
		CreateJson(w, "Please enter at least one lowercase letter", http.StatusBadRequest)
		return false
	}
	reg = regexp.MustCompile(`[A-Z]`)
	if !reg.MatchString(password) {
		CreateJson(w, "Please enter at least one uppercase letter", http.StatusBadRequest)
		return false
	}

	reg = regexp.MustCompile(`[\d]`)
	if !reg.MatchString(password) {
		CreateJson(w, "Please enter at least one digit", http.StatusBadRequest)
		return false
	}

	reg = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	if !reg.MatchString(password) {
		CreateJson(w, "Please enter at least one special character", http.StatusBadRequest)
		return false

	}

	return true
}

func IsValidCheckPassword(w http.ResponseWriter, pass1 string, pass2 string) bool {
	if pass2 == "" {
		CreateJson(w, "Please enter your password", http.StatusBadRequest)
		return false
	}
	if pass1 != pass2 {
		CreateJson(w, "your password not correct", http.StatusBadRequest)
		return false
	}
	return true
}
