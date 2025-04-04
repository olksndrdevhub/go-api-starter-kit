package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
)

type passwordConfig struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var config = &passwordConfig{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

type PasswordValidationError struct {
	Message string
}

func (pve PasswordValidationError) Error() string {
	return pve.Message
}

// ValidatePassword validates the password by length, complexity, etc
func ValidatePassword(password string) error {
	// check length
	if len(password) < 8 {
		return &PasswordValidationError{
			Message: "Password must be at least 8 characters long",
		}
	}
	if len(password) > 72 {
		return &PasswordValidationError{
			Message: "Password must be less than 72 characters long",
		}
	}

	var (
		hasUpper,
		hasLower,
		hasNumber,
		hasSpecial bool
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
		return &PasswordValidationError{
			Message: "Password must contain at least one uppercase letter",
		}
	}
	if !hasLower {
		return &PasswordValidationError{
			Message: "Password must contain at least one lowercase letter",
		}
	}
	if !hasNumber {
		return &PasswordValidationError{
			Message: "Password must contain at least one number",
		}
	}
	if !hasSpecial {
		return &PasswordValidationError{
			Message: "Password must contain at least one special character",
		}
	}
	return nil
}

// HashPassword generates a secure hash from the password
func HashPassword(password string) (string, error) {
	// random salt
	salt := make([]byte, config.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// hash password using argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		config.iterations,
		config.memory,
		config.parallelism,
		config.keyLength,
	)

	// format the hash with its params for storage
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.memory,
		config.iterations,
		config.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encodedHash, nil

}

func VerifyPassword(password, encodedHash string) (bool, error) {
	// extract the parameters from the encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	if version != argon2.Version {
		return false, fmt.Errorf("invalid hash version")
	}

	// parse memory iterations and parallelism
	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	// decode salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	keyLength := uint32(len(storedHash))

	// compute hash from provided password with same parameters
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	return subtle.ConstantTimeCompare(storedHash, computedHash) == 1, nil

}
