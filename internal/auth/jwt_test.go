package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test JWT config
var testConfig = &JWTConfig{
	AccessTokenSecret:  "test-access-secret-key-12345",
	RefreshTokenSecret: "test-refresh-secret-key-67890",
	AccessTokenExpiry:  15 * time.Minute,
	RefreshTokenExpiry: 7 * 24 * time.Hour,
}

func TestJWTConfig_GenerateAccessToken(t *testing.T) {
	t.Run("should generate valid access token", func(t *testing.T) {
		userID := "user-123"
		userType := "trainer"
		email := "trainer@fitterby.com"

		token, err := testConfig.GenerateAccessToken(userID, userType, email)

		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Verify the token can be validated
		claims, valid, err := testConfig.ValidateAccessToken(token)
		require.NoError(t, err)
		require.True(t, valid)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, userType, claims.UserType)
		assert.Equal(t, email, claims.Email)
		assert.WithinDuration(t, time.Now().Add(testConfig.AccessTokenExpiry), claims.ExpiresAt.Time, time.Second)
	})

	t.Run("should fail with empty user ID", func(t *testing.T) {
		token, err := testConfig.GenerateAccessToken("", "trainer", "test@fitterby.com")

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("should fail with empty user type", func(t *testing.T) {
		token, err := testConfig.GenerateAccessToken("user-123", "", "test@fitterby.com")

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("should fail with empty email", func(t *testing.T) {
		token, err := testConfig.GenerateAccessToken("user-123", "trainer", "")

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestJWTConfig_GenerateRefreshToken(t *testing.T) {
	t.Run("should generate valid refresh token", func(t *testing.T) {
		userID := "user-456"
		userType := "client"
		email := "client@fitterby.com"

		token, err := testConfig.GenerateRefreshToken(userID, userType, email)

		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Verify the token can be validated
		claims, valid, err := testConfig.ValidateRefreshToken(token)
		require.NoError(t, err)
		require.True(t, valid)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, userType, claims.UserType)
		assert.Equal(t, email, claims.Email)
		assert.WithinDuration(t, time.Now().Add(testConfig.RefreshTokenExpiry), claims.ExpiresAt.Time, time.Second)
	})

	t.Run("should have different secrets for access and refresh tokens", func(t *testing.T) {
		userID := "user-123"
		userType := "trainer"
		email := "test@fitterby.com"

		accessToken, err := testConfig.GenerateAccessToken(userID, userType, email)
		require.NoError(t, err)

		refreshToken, err := testConfig.GenerateRefreshToken(userID, userType, email)
		require.NoError(t, err)

		// Access token should NOT validate as refresh token
		_, valid, _ := testConfig.ValidateRefreshToken(accessToken)
		assert.False(t, valid)

		// Refresh token should NOT validate as access token
		_, valid, _ = testConfig.ValidateAccessToken(refreshToken)
		assert.False(t, valid)
	})
}

func TestJWTConfig_ValidateAccessToken(t *testing.T) {
	t.Run("should validate valid access token", func(t *testing.T) {
		userID := "user-789"
		userType := "client"
		email := "client@fitterby.com"

		token, err := testConfig.GenerateAccessToken(userID, userType, email)
		require.NoError(t, err)

		claims, valid, err := testConfig.ValidateAccessToken(token)

		require.NoError(t, err)
		assert.True(t, valid)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, userType, claims.UserType)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("should reject expired access token", func(t *testing.T) {
		// Create an expired token manually
		expiredClaims := &Claims{
			UserID:   "user-123",
			UserType: "trainer",
			Email:    "trainer@fitterby.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		tokenString, err := token.SignedString([]byte(testConfig.AccessTokenSecret))
		require.NoError(t, err)

		claims, valid, err := testConfig.ValidateAccessToken(tokenString)

		assert.Error(t, err)
		assert.False(t, valid)
		assert.Nil(t, claims)
	})

	t.Run("should reject tampered access token", func(t *testing.T) {
		validToken, err := testConfig.GenerateAccessToken("user-123", "client", "test@fitterby.com")
		require.NoError(t, err)

		// Tamper with the token
		tamperedToken := validToken[:len(validToken)-1] + "x"

		claims, valid, err := testConfig.ValidateAccessToken(tamperedToken)

		assert.Error(t, err)
		assert.False(t, valid)
		assert.Nil(t, claims)
	})

	t.Run("should reject token with wrong secret", func(t *testing.T) {
		// Create token with wrong secret
		wrongConfig := &JWTConfig{
			AccessTokenSecret:  "wrong-secret-key",
			RefreshTokenSecret: testConfig.RefreshTokenSecret,
			AccessTokenExpiry:  testConfig.AccessTokenExpiry,
			RefreshTokenExpiry: testConfig.RefreshTokenExpiry,
		}

		token, err := wrongConfig.GenerateAccessToken("user-123", "trainer", "test@fitterby.com")
		require.NoError(t, err)

		// Try to validate with correct config (should fail)
		claims, valid, err := testConfig.ValidateAccessToken(token)

		assert.Error(t, err)
		assert.False(t, valid)
		assert.Nil(t, claims)
	})

	t.Run("should reject malformed token", func(t *testing.T) {
		claims, valid, err := testConfig.ValidateAccessToken("not-a-valid-token")

		assert.Error(t, err)
		assert.False(t, valid)
		assert.Nil(t, claims)
	})
}

func TestJWTConfig_ValidateRefreshToken(t *testing.T) {
	t.Run("should validate valid refresh token", func(t *testing.T) {
		userID := "user-999"
		userType := "trainer"
		email := "trainer@fitterby.com"

		token, err := testConfig.GenerateRefreshToken(userID, userType, email)
		require.NoError(t, err)

		claims, valid, err := testConfig.ValidateRefreshToken(token)

		require.NoError(t, err)
		assert.True(t, valid)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, userType, claims.UserType)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("should reject access token as refresh token", func(t *testing.T) {
		accessToken, err := testConfig.GenerateAccessToken("user-123", "client", "test@fitterby.com")
		require.NoError(t, err)

		claims, valid, err := testConfig.ValidateRefreshToken(accessToken)

		assert.Error(t, err)
		assert.False(t, valid)
		assert.Nil(t, claims)
	})
}

func TestJWTConfig_RefreshTokens(t *testing.T) {
	t.Run("should refresh valid tokens", func(t *testing.T) {
		userID := "user-123"
		userType := "client"
		email := "client@fitterby.com"

		// Generate initial refresh token
		initialRefreshToken, err := testConfig.GenerateRefreshToken(userID, userType, email)
		require.NoError(t, err)

		// Refresh tokens
		newAccessToken, newRefreshToken, err := testConfig.RefreshTokens(initialRefreshToken)

		require.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
		assert.NotEqual(t, initialRefreshToken, newRefreshToken)

		// Verify new tokens work
		accessClaims, valid, err := testConfig.ValidateAccessToken(newAccessToken)
		require.NoError(t, err)
		assert.True(t, valid)
		assert.Equal(t, userID, accessClaims.UserID)

		refreshClaims, valid, err := testConfig.ValidateRefreshToken(newRefreshToken)
		require.NoError(t, err)
		assert.True(t, valid)
		assert.Equal(t, userID, refreshClaims.UserID)
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		newAccessToken, newRefreshToken, err := testConfig.RefreshTokens("invalid-token")

		assert.Error(t, err)
		assert.Empty(t, newAccessToken)
		assert.Empty(t, newRefreshToken)
	})

	t.Run("should fail with expired refresh token", func(t *testing.T) {
		// Create an expired refresh token
		expiredClaims := &Claims{
			UserID:   "user-123",
			UserType: "trainer",
			Email:    "trainer@fitterby.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredToken, err := token.SignedString([]byte(testConfig.RefreshTokenSecret))
		require.NoError(t, err)

		newAccessToken, newRefreshToken, err := testConfig.RefreshTokens(expiredToken)

		assert.Error(t, err)
		assert.Empty(t, newAccessToken)
		assert.Empty(t, newRefreshToken)
	})

	t.Run("should fail with access token instead of refresh token", func(t *testing.T) {
		accessToken, err := testConfig.GenerateAccessToken("user-123", "client", "test@fitterby.com")
		require.NoError(t, err)

		newAccessToken, newRefreshToken, err := testConfig.RefreshTokens(accessToken)

		assert.Error(t, err)
		assert.Empty(t, newAccessToken)
		assert.Empty(t, newRefreshToken)
	})
}

func TestJWTConfig_TokenExpiry(t *testing.T) {
	t.Run("access token should have correct expiry", func(t *testing.T) {
		userID := "user-123"
		userType := "trainer"
		email := "trainer@fitterby.com"

		startTime := time.Now()
		token, err := testConfig.GenerateAccessToken(userID, userType, email)
		require.NoError(t, err)

		claims, valid, err := testConfig.ValidateAccessToken(token)
		require.NoError(t, err)
		require.True(t, valid)

		expectedExpiry := startTime.Add(testConfig.AccessTokenExpiry)
		actualExpiry := claims.ExpiresAt.Time

		assert.WithinDuration(t, expectedExpiry, actualExpiry, time.Second)
	})

	t.Run("refresh token should have correct expiry", func(t *testing.T) {
		userID := "user-456"
		userType := "client"
		email := "client@fitterby.com"

		startTime := time.Now()
		token, err := testConfig.GenerateRefreshToken(userID, userType, email)
		require.NoError(t, err)

		claims, valid, err := testConfig.ValidateRefreshToken(token)
		require.NoError(t, err)
		require.True(t, valid)

		expectedExpiry := startTime.Add(testConfig.RefreshTokenExpiry)
		actualExpiry := claims.ExpiresAt.Time

		assert.WithinDuration(t, expectedExpiry, actualExpiry, time.Second)
	})
}

// Edge case tests
func TestJWTConfig_EdgeCases(t *testing.T) {
	t.Run("should handle very short token expiry", func(t *testing.T) {
		shortConfig := &JWTConfig{
			AccessTokenSecret:  "short-secret",
			RefreshTokenSecret: "short-refresh-secret",
			AccessTokenExpiry:  1 * time.Second,
			RefreshTokenExpiry: 2 * time.Second,
		}

		token, err := shortConfig.GenerateAccessToken("user-123", "trainer", "test@fitterby.com")
		require.NoError(t, err)

		claims, valid, err := shortConfig.ValidateAccessToken(token)
		require.NoError(t, err)
		assert.True(t, valid)
		assert.Equal(t, "user-123", claims.UserID)
	})

	t.Run("should handle different config instances", func(t *testing.T) {
		config1 := &JWTConfig{
			AccessTokenSecret:  "config1-secret",
			RefreshTokenSecret: "config1-refresh",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}

		config2 := &JWTConfig{
			AccessTokenSecret:  "config2-secret",
			RefreshTokenSecret: "config2-refresh",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}

		// Generate token with config1
		token, err := config1.GenerateAccessToken("user-123", "trainer", "test@fitterby.com")
		require.NoError(t, err)

		// Should not validate with config2
		claims, valid, err := config2.ValidateAccessToken(token)
		assert.Error(t, err)
		assert.False(t, valid)
		assert.Nil(t, claims)
	})
}

// Benchmark tests
func BenchmarkGenerateAccessToken(b *testing.B) {
	userID := "user-123"
	userType := "trainer"
	email := "trainer@fitterby.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := testConfig.GenerateAccessToken(userID, userType, email)
		if err != nil {
			b.Fatalf("Failed to generate token: %v", err)
		}
	}
}

func BenchmarkValidateAccessToken(b *testing.B) {
	userID := "user-123"
	userType := "trainer"
	email := "trainer@fitterby.com"

	token, err := testConfig.GenerateAccessToken(userID, userType, email)
	if err != nil {
		b.Fatalf("Failed to generate token: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := testConfig.ValidateAccessToken(token)
		if err != nil {
			b.Fatalf("Failed to validate token: %v", err)
		}
	}
}
