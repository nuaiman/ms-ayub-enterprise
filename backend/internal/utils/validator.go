package utils

import "regexp"

func IsValidEmail(email string) bool {
	// Basic email validation regex
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func IsValidPhone(phone string) bool {
	// Basic phone validation - at least 7 digits, only numbers, +, -, spaces
	pattern := `^[\+\d\s\-\(\)]{7,20}$`
	match, _ := regexp.MatchString(pattern, phone)
	return match
}

func IsValidIDType(idType string) bool {
	validTypes := map[string]bool{
		"nid":               true,
		"passport":          true,
		"driving_license":   true,
		"birth_certificate": true,
		"trade_license":     true,
		"other":             true,
	}
	return validTypes[idType]
}

func IsValidRole(role string) bool {
	validRoles := map[string]bool{
		"admin": true, "manager": true, "accounts": true, "staff": true,
	}
	return validRoles[role]
}
