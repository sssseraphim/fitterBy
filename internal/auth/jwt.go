package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func (cfg *JWTConfig) GenerateAccessToken(userId, userType, email string) (string, error) {
	if userId == "" || userType == "" || email == "" {
		return "", errors.New("token requires user id, type and email")
	}
	expirationTime := time.Now().Add(cfg.AccessTokenExpiry)
	claims := Claims{
		UserID:   userId,
		UserType: userType,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.AccessTokenSecret))
}

func (cfg *JWTConfig) GenerateRefreshToken(userId, userType, email string) (string, error) {
	if userId == "" || userType == "" || email == "" {
		return "", errors.New("token requires user id, type and email")
	}
	expirationTime := time.Now().Add(cfg.RefreshTokenExpiry)
	claims := Claims{
		UserID:   userId,
		UserType: userType,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.RefreshTokenSecret))
}

func (cfg *JWTConfig) ValidateAccessToken(tokenString string) (claims *Claims, valid bool, err error) {
	claims = &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.AccessTokenSecret), nil
	})
	if err != nil {
		return nil, false, err
	}
	if !token.Valid {
		return nil, false, err
	}
	return claims, true, nil
}

func (cfg *JWTConfig) ValidateRefreshToken(tokenString string) (claims *Claims, valid bool, err error) {
	claims = &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.RefreshTokenSecret), nil
	})
	if err != nil {
		return nil, false, err
	}
	if !token.Valid {
		return nil, false, err
	}
	return claims, true, nil
}

func (cfg *JWTConfig) RefreshTokens(refreshToken string) (string, string, error) {
	claims, valid, err := cfg.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	if !valid {
		return "", "", nil
	}
	newAccessToken, err := cfg.GenerateAccessToken(claims.UserID, claims.UserType, claims.Email)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := cfg.GenerateRefreshToken(claims.UserID, claims.UserType, claims.Email)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil

}
