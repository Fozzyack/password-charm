// Package utils provides utility functions for the password manager.
// This includes password generation, filename creation, and other helper functions.
package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

// PasswordOptions configures password generation parameters
type PasswordOptions struct {
	Length            int  // Password length (8-64)
	IncludeUppercase  bool // Include A-Z
	IncludeLowercase  bool // Include a-z
	IncludeNumbers    bool // Include 0-9
	IncludeSymbols    bool // Include special symbols
	ExcludeAmbiguous  bool // Exclude ambiguous characters like 0, O, l, I
}

// Character sets for password generation
const (
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	numberChars    = "0123456789"
	symbolChars    = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	ambiguousChars = "0O1lI"
)

// DefaultPasswordOptions returns sensible default options for password generation
func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{
		Length:            16,
		IncludeUppercase:  true,
		IncludeLowercase:  true,
		IncludeNumbers:    true,
		IncludeSymbols:    true,
		ExcludeAmbiguous:  true,
	}
}

// GeneratePassword creates a secure random password based on the given options
func GeneratePassword(opts PasswordOptions) (string, error) {
	// Validate options
	if opts.Length < 8 || opts.Length > 64 {
		return "", fmt.Errorf("password length must be between 8 and 64 characters")
	}

	if !opts.IncludeUppercase && !opts.IncludeLowercase && !opts.IncludeNumbers && !opts.IncludeSymbols {
		return "", fmt.Errorf("at least one character type must be included")
	}

	// Build character set
	var charset string
	var guaranteedChars string

	if opts.IncludeUppercase {
		chars := uppercaseChars
		if opts.ExcludeAmbiguous {
			chars = removeAmbiguousChars(chars)
		}
		charset += chars
		// Guarantee at least one uppercase character
		if char, err := getRandomChar(chars); err == nil {
			guaranteedChars += char
		}
	}

	if opts.IncludeLowercase {
		chars := lowercaseChars
		if opts.ExcludeAmbiguous {
			chars = removeAmbiguousChars(chars)
		}
		charset += chars
		// Guarantee at least one lowercase character
		if char, err := getRandomChar(chars); err == nil {
			guaranteedChars += char
		}
	}

	if opts.IncludeNumbers {
		chars := numberChars
		if opts.ExcludeAmbiguous {
			chars = removeAmbiguousChars(chars)
		}
		charset += chars
		// Guarantee at least one number
		if char, err := getRandomChar(chars); err == nil {
			guaranteedChars += char
		}
	}

	if opts.IncludeSymbols {
		charset += symbolChars
		// Guarantee at least one symbol
		if char, err := getRandomChar(symbolChars); err == nil {
			guaranteedChars += char
		}
	}

	// Generate password
	password := make([]byte, opts.Length)
	
	// First, place guaranteed characters
	for i, char := range guaranteedChars {
		if i < len(password) {
			password[i] = byte(char)
		}
	}

	// Fill remaining positions with random characters
	for i := len(guaranteedChars); i < opts.Length; i++ {
		randomChar, err := getRandomChar(charset)
		if err != nil {
			return "", fmt.Errorf("failed to generate random character: %v", err)
		}
		password[i] = randomChar[0]
	}

	// Shuffle the password to randomize guaranteed character positions
	if err := shuffleBytes(password); err != nil {
		return "", fmt.Errorf("failed to shuffle password: %v", err)
	}

	return string(password), nil
}

// removeAmbiguousChars removes potentially confusing characters from a charset
func removeAmbiguousChars(charset string) string {
	result := charset
	for _, ambiguous := range ambiguousChars {
		result = strings.ReplaceAll(result, string(ambiguous), "")
	}
	return result
}

// getRandomChar returns a random character from the given charset
func getRandomChar(charset string) (string, error) {
	if len(charset) == 0 {
		return "", fmt.Errorf("charset is empty")
	}

	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return "", err
	}

	return string(charset[randomIndex.Int64()]), nil
}

// shuffleBytes randomly shuffles a byte slice using Fisher-Yates algorithm
func shuffleBytes(slice []byte) error {
	for i := len(slice) - 1; i > 0; i-- {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		j := randomIndex.Int64()
		slice[i], slice[j] = slice[j], slice[i]
	}
	return nil
}

// EvaluatePasswordStrength returns a strength score (0-4) and description for a password
func EvaluatePasswordStrength(password string) (int, string) {
	score := 0
	
	// Length check
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}

	// Character variety checks
	if matched, _ := regexp.MatchString(`[a-z]`, password); matched {
		score++
	}
	if matched, _ := regexp.MatchString(`[A-Z]`, password); matched {
		score++
	}
	if matched, _ := regexp.MatchString(`[0-9]`, password); matched {
		score++
	}
	if matched, _ := regexp.MatchString(`[^a-zA-Z0-9]`, password); matched {
		score++
	}

	// Normalize score to 0-4 range
	if score > 4 {
		score = 4
	}

	descriptions := []string{
		"Very Weak",
		"Weak", 
		"Fair",
		"Good",
		"Strong",
	}

	return score, descriptions[score]
}

// GenerateFilename creates a unique filename for storing password entries
func GenerateFilename(siteName string) string {
	// Clean the site name for use as filename
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	cleanName := reg.ReplaceAllString(siteName, "_")
	cleanName = strings.ToLower(cleanName)
	
	// Limit length
	if len(cleanName) > 20 {
		cleanName = cleanName[:20]
	}
	
	// Add timestamp for uniqueness
	timestamp := time.Now().Format("20060102_150405")
	
	// Combine name and timestamp
	filename := fmt.Sprintf("%s_%s", cleanName, timestamp)
	
	return filename
}

// SanitizeInput removes potentially dangerous characters from user input
func SanitizeInput(input string) string {
	// Remove control characters and normalize whitespace
	reg := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	sanitized := reg.ReplaceAllString(input, "")
	
	// Normalize whitespace
	sanitized = strings.TrimSpace(sanitized)
	
	return sanitized
}

// ParseFilenameToSiteName converts a filename back to a readable site name
func ParseFilenameToSiteName(filename string) string {
	// Remove .gpg extension if present
	if strings.HasSuffix(filename, ".gpg") {
		filename = strings.TrimSuffix(filename, ".gpg")
	}
	
	// Find the last underscore (timestamp separator)
	lastUnderscore := strings.LastIndex(filename, "_")
	if lastUnderscore == -1 {
		// No timestamp found, return as is with underscores replaced by spaces
		return strings.ReplaceAll(filename, "_", " ")
	}
	
	// Check if what follows the last underscore looks like a timestamp
	potentialTimestamp := filename[lastUnderscore+1:]
	if len(potentialTimestamp) == 15 && strings.Contains(potentialTimestamp, "_") {
		// This looks like our timestamp format (YYYYMMDD_HHMMSS)
		siteName := filename[:lastUnderscore]
		return strings.ReplaceAll(siteName, "_", " ")
	}
	
	// No valid timestamp found, return as is with underscores replaced
	return strings.ReplaceAll(filename, "_", " ")
}

// FormatTimestampForDisplay formats a time.Time for user-friendly display
func FormatTimestampForDisplay(t time.Time) string {
	now := time.Now()
	
	// If it's today, show time
	if t.Format("2006-01-02") == now.Format("2006-01-02") {
		return "Today " + t.Format("3:04 PM")
	}
	
	// If it's yesterday, show "Yesterday"
	yesterday := now.AddDate(0, 0, -1)
	if t.Format("2006-01-02") == yesterday.Format("2006-01-02") {
		return "Yesterday " + t.Format("3:04 PM")
	}
	
	// If it's within a week, show day name
	if now.Sub(t).Hours() < 7*24 {
		return t.Format("Monday 3:04 PM")
	}
	
	// If it's within this year, show month and day
	if t.Year() == now.Year() {
		return t.Format("Jan 2, 3:04 PM")
	}
	
	// Otherwise show full date
	return t.Format("Jan 2, 2006")
}

// TruncateString truncates a string to the specified length and adds ellipsis if needed
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}