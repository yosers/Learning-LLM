package jwt

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID int32    `json:"user_id"`
	Role   []string `json:"role"`
	jwt.RegisteredClaims
}

// TokenBlacklist stores invalidated tokens
var (
	blacklistedTokens = make(map[string]time.Time)
	blacklistMutex    sync.RWMutex
)

// Clean up expired tokens from blacklist periodically
func init() {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			cleanupBlacklist()
		}
	}()
}

func cleanupBlacklist() {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	now := time.Now()
	for token, expiry := range blacklistedTokens {
		if now.After(expiry) {
			delete(blacklistedTokens, token)
		}
	}
}

func GenerateToken(userID int32, role []string) (string, error) {
	// Get secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	// Create claims with user ID and standard claims
	claims := JWTClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*JWTClaim, error) {
	// Check if token is blacklisted
	blacklistMutex.RLock()
	if _, blacklisted := blacklistedTokens[tokenString]; blacklisted {
		blacklistMutex.RUnlock()
		return nil, fmt.Errorf("token has been invalidated")
	}
	blacklistMutex.RUnlock()

	// Get secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// InvalidateToken adds a token to the blacklist
func InvalidateToken(tokenString string) error {
	claims := &JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		secretKey := os.Getenv("JWT_SECRET_KEY")
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token")
	}

	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	// Store token in blacklist with its expiry time
	if claims.ExpiresAt != nil {
		blacklistedTokens[tokenString] = claims.ExpiresAt.Time
	} else {
		// If no expiry, blacklist for 24 hours
		blacklistedTokens[tokenString] = time.Now().Add(24 * time.Hour)
	}

	return nil
}
