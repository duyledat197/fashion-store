package crypto_util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Constants for generating random password
const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes     = "0123456789"
)

// GeneratePassword generates a random password based on the specified criteria
func GeneratePassword(length int, useLetters bool, useSpecial bool, useNum bool) string {
	b := make([]byte, length)
	for i := range b {
		if useLetters {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		} else if useSpecial {
			b[i] = specialBytes[rand.Intn(len(specialBytes))]
		} else if useNum {
			b[i] = numBytes[rand.Intn(len(numBytes))]
		}
	}
	return string(b)
}

// Constants for generating a code
const (
	Charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	RandLength = 5
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateCode generates a unique code with a prefix based on the current timestamp
func GenerateCode(prefix string) string {
	t := time.Now()
	y := ""
	if t.Year()%100 < 10 {
		y = fmt.Sprintf("0%d", t.Year()%100)
	} else {
		y = fmt.Sprintf("%d", t.Year()%100)
	}
	m := ""
	if t.Month() < 10 {
		m = fmt.Sprintf("0%d", t.Month())
	} else {
		m = fmt.Sprintf("%d", t.Month())
	}
	d := ""
	if t.Day() < 10 {
		d = fmt.Sprintf("0%d", t.Day())
	} else {
		d = fmt.Sprintf("%d", t.Day())
	}
	code := fmt.Sprintf("%s%s%s%s%s", prefix, y, m, d, genStringWithLength(RandLength))
	return strings.ToUpper(code)
}

// Utility function to generate a string with characters from a predefined charset
func stringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[seededRand.Intn(len(Charset))]
	}
	return string(b)
}

// Utility function to generate a string with a specified length
func genStringWithLength(length int) string {
	return stringWithCharset(length)
}
