package utils

import (
	"errors"
	"time"

	"github.com/devjoemedia/go-ticketing-api/config"
	"github.com/golang-jwt/jwt/v5"
)

// Custom Claims
type AccessClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID  uint   `json:"user_id"`
	TokenID string `json:"token_id"`
	jwt.RegisteredClaims
}

// Helpers
func accessSecret() []byte     { return []byte(config.AppConfig.JWTAccessSecret) }
func refreshSecret() []byte    { return []byte(config.AppConfig.JWTRefreshSecret) }
func accessExpiryMinutes() int { return int(config.AppConfig.JWTAccessExpiryMinutes) }
func refreshExpiryDays() int   { return int(config.AppConfig.JWTRefreshExpiryDays) }

func accessExpiry() time.Duration {
	mins := accessExpiryMinutes()
	if mins == 0 {
		mins = 15
	}
	return time.Duration(mins) * time.Minute
}

func refreshExpiry() time.Duration {
	days := refreshExpiryDays()
	if days == 0 {
		days = 7
	}
	return time.Duration(days) * 24 * time.Hour
}

// Token Generation
func GenerateAccessToken(userID uint, email string) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiry())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chitodo-api",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(accessSecret())
}

func GenerateRefreshToken(userID uint, tokenID string) (string, error) {
	claims := RefreshClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chitodo-api",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(refreshSecret())
}

// Token Parsing
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

func ParseAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return accessSecret(), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return refreshSecret(), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// RefreshExpiry exported for handler use
func RefreshExpiry() time.Duration { return refreshExpiry() }
