package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sssseraphim/fitterBy/internal/auth"
	"github.com/sssseraphim/fitterBy/internal/database"
)

type TokenService struct {
	DB        *database.Queries
	JWTConfig *auth.JWTConfig
}

func NewTokenService(db *database.Queries, jwtConfig *auth.JWTConfig) *TokenService {
	return &TokenService{
		DB:        db,
		JWTConfig: jwtConfig}
}

func (s *TokenService) GenerateTokens(ctx context.Context, userID, userType, email string) (string, string, error) {
	accessToken, err := s.JWTConfig.GenerateAccessToken(userID, userType, email)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.JWTConfig.GenerateRefreshToken(userID, userType, email)
	if err != nil {
		return "", "", err
	}
	hash := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := hex.EncodeToString(hash[:])
	expires_at := time.Now().Add(s.JWTConfig.RefreshTokenExpiry)
	_, err = s.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID:    uuid.MustParse(userID),
		TokenHash: refreshTokenHash,
		ExpiresAt: expires_at,
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) ValidateAndRefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	hash := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := hex.EncodeToString(hash[:])

	storedToken, err := s.DB.GetRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		return "", "", err
	}

	if time.Now().After(storedToken.ExpiresAt) {
		s.DB.DeleteRefreshToken(ctx, storedToken.TokenHash)
		return "", "", errors.New("token expired")
	}

	claims, valid, err := s.JWTConfig.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	if !valid {
		return "", "", errors.New("Invalid Token")
	}

	newAccessToken, newRefreshToken, err := s.GenerateTokens(ctx, claims.UserID, claims.UserType, claims.Email)
	if err != nil {
		return "", "", err
	}
	s.DB.DeleteRefreshToken(ctx, storedToken.TokenHash)
	return newAccessToken, newRefreshToken, nil

}

func (s *TokenService) RevokeUserTokens(ctx context.Context, userID string) error {
	return s.DB.DeleteUserRefreshToken(ctx, uuid.MustParse(userID))
}

func (s *TokenService) CleanExpiredTokens(ctx context.Context) error {
	return s.DB.CleanExpiredRefreshTokens(ctx)
}
