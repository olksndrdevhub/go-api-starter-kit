package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var secretKey = []byte("secretkey") // set secure secret key in prod

// JWTClaims represents the claims in a JWT
type JWTClaims struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"expires_at"` // in seconds
}

func GenerateJWTToken(userID int64, email string) (string, error) {
	// create header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)

	// create payload with claims
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // expire in 24 hours
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// create signature
	signatureInput := encodedHeader + "." + encodedPayload
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(signatureInput))
	signature := h.Sum(nil)
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	// combine to form JWT token
	token := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, encodedSignature)
	return token, nil
}

// ValidateJWTToken validates a JWT token and returns the claims if valid
func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	encodedHeader, encodedPayload, encodedSignature := parts[0], parts[1], parts[2]

	// verify signature
	signatureInput := encodedHeader + "." + encodedPayload
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(signatureInput))
	signature := h.Sum(nil)
	expectedSignature := base64.RawURLEncoding.EncodeToString(signature)

	if encodedSignature != expectedSignature {
		return nil, errors.New("invalid token signature")
	}

	// decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, err
	}

	var claims JWTClaims

	err = json.Unmarshal(payloadJSON, &claims)
	if err != nil {
		return nil, err
	}

	// check expiration time
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}
