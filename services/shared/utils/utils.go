package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// GenerateID generates a new UUID
func GenerateID() string {
	return uuid.New().String()
}

// GenerateShortID generates a shorter ID (8 characters)
func GenerateShortID() string {
	return uuid.New().String()[:8]
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}

	return string(bytes), nil
}

// HashPassword creates a SHA256 hash of a password (in production, use bcrypt)
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidatePhoneNumber validates a phone number format (basic validation)
func ValidatePhoneNumber(phone string) bool {
	// Remove all non-digit characters
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Check if it's between 10-15 digits
	return len(cleaned) >= 10 && len(cleaned) <= 15
}

// GenerateSKU generates a product SKU
func GenerateSKU(prefix string) string {
	timestamp := time.Now().Unix()
	random, _ := GenerateRandomString(4)
	return fmt.Sprintf("%s-%d-%s", strings.ToUpper(prefix), timestamp, strings.ToUpper(random))
}

// SlugifyString converts a string to a URL-friendly slug
func SlugifyString(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	text = reg.ReplaceAllString(text, "-")

	// Remove leading and trailing hyphens
	text = strings.Trim(text, "-")

	return text
}

// TruncateString truncates a string to specified length
func TruncateString(text string, length int) string {
	if len(text) <= length {
		return text
	}

	if length <= 3 {
		return text[:length]
	}

	return text[:length-3] + "..."
}

// StringInSlice checks if a string exists in a slice
func StringInSlice(str string, slice []string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// RemoveStringFromSlice removes a string from a slice
func RemoveStringFromSlice(slice []string, str string) []string {
	result := make([]string, 0)
	for _, item := range slice {
		if item != str {
			result = append(result, item)
		}
	}
	return result
}

// UniqueStrings returns unique strings from a slice
func UniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	result := make([]string, 0)

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// FormatCurrency formats a float64 as currency string
func FormatCurrency(amount float64, currency string) string {
	switch currency {
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	case "EUR":
		return fmt.Sprintf("€%.2f", amount)
	case "GBP":
		return fmt.Sprintf("£%.2f", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, currency)
	}
}

// ParseCurrencyAmount extracts numeric amount from currency string
func ParseCurrencyAmount(currency string) (float64, error) {
	// Remove currency symbols and spaces
	cleaned := regexp.MustCompile(`[^0-9.]`).ReplaceAllString(currency, "")

	var amount float64
	_, err := fmt.Sscanf(cleaned, "%f", &amount)
	return amount, err
}

// CalculatePercentage calculates percentage between two numbers
func CalculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (part / total) * 100
}

// RoundToDecimal rounds a float64 to specified decimal places
func RoundToDecimal(value float64, decimals int) float64 {
	multiplier := float64(1)
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	return float64(int(value*multiplier+0.5)) / multiplier
}

// TimePointer returns a pointer to time.Time
func TimePointer(t time.Time) *time.Time {
	return &t
}

// StringPointer returns a pointer to string
func StringPointer(s string) *string {
	return &s
}

// IntPointer returns a pointer to int
func IntPointer(i int) *int {
	return &i
}

// BoolPointer returns a pointer to bool
func BoolPointer(b bool) *bool {
	return &b
}

// DerefString safely dereferences a string pointer
func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// DerefInt safely dereferences an int pointer
func DerefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// DerefBool safely dereferences a bool pointer
func DerefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// DerefTime safely dereferences a time pointer
func DerefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}

// MaskEmail masks an email address for privacy
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return email
	}

	masked := username[:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:]
	return masked + "@" + domain
}

// MaskCreditCard masks a credit card number
func MaskCreditCard(cardNumber string) string {
	if len(cardNumber) < 4 {
		return cardNumber
	}

	return strings.Repeat("*", len(cardNumber)-4) + cardNumber[len(cardNumber)-4:]
}
