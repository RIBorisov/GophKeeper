package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/RIBorisov/GophKeeper/internal/log"
)

// Claims represents the claims for a JWT token.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// GetUserID retrieves a userID string from the passed JWT token.
func (s *Service) GetUserID(tokenString, secretKey string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Error("failed parse with claims tokenString", "err", err)
		return ""
	}
	if !token.Valid {
		log.Error("Invalid token", "token", token)
		return ""
	}

	return claims.UserID
}

// BuildJWTString generates a JWT token.
func BuildJWTString(secretKey, userID string) (string, error) {
	const tokenExp = time.Hour * 720

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create token string: %w", err)
	}

	return tokenString, nil
}

func hashPassword(secret, data string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data+secret), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed hash password: %w", err)
	}

	return string(hashed), nil
}

func comparePasswords(secret, hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+secret))
	if err != nil {
		return ErrInvalidPassword
	}

	return nil
}
