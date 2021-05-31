package utils

import (
	"regexp"
	"strconv"
	"unicode"
)

var (
	phoneNumberRegex = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	emailRegex       = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func IsValidPhoneNumber(number string) bool {

	return phoneNumberRegex.MatchString(number)
}

func IsValidEmial(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func IsValidPassword(field, s string) (isValid bool, errs []string) {
	var (
		isMin   bool
		special bool
		number  bool
		upper   bool
		lower   bool
	)
	min := 6
	max := 20

	// append error
	appendError := func(err string) {
		errs = append(errs, field+" "+err)
	}

	if len(s) < min || len(s) > max {
		isMin = false
		appendError("length should be " + strconv.Itoa(min) + " to " + strconv.Itoa(max))
	}

	for _, c := range s {
		// Optimize perf if all become true before reaching the end
		if special && number && upper && lower && isMin {
			break
		}

		// else go on switching
		switch {
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsNumber(c):
			number = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}

	// Add custom error messages
	if !special {
		appendError("should contain at least a single special character")
	}
	if !number {
		appendError("should contain at least a single digit")
	}
	if !lower {
		appendError("should contain at least a single lowercase letter")
	}
	if !upper {
		appendError("should contain at least single uppercase letter")
	}

	// if there is any error
	if len(errs) > 0 {
		return false, errs
	}

	// everyting is right
	return true, errs
}
