package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/sssseraphim/fitterBy/internal/auth"
	"github.com/sssseraphim/fitterBy/internal/database"
	"github.com/sssseraphim/fitterBy/internal/middleware"
	"github.com/sssseraphim/fitterBy/internal/services"
)

type AuthHandler struct {
	DB           *database.Queries
	JWTConfig    *auth.JWTConfig
	TokenService services.TokenService
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	TokenResoponse *TokenResoponse
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResoponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, TokenResoponse{})
		return
	}
	if req.RefreshToken == "" {
		respondWithJSON(w, http.StatusBadRequest, TokenResoponse{})
		return
	}
	accessToken, refreshToken, err := h.TokenService.ValidateAndRefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, TokenResoponse{})
		return
	}
	response := TokenResoponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.TokenService.JWTConfig.AccessTokenExpiry.Seconds()),
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)

	if err := h.TokenService.RevokeUserTokens(r.Context(), userId); err != nil {
		respondWithError(w, 500, "failed to logout", err)
		return
	}
	respondWithJSON(w, 200, map[string]string{"message": "Logged out successfully"})
}

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, _, err := argon2id.CheckHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "bad request", err)
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "All fields are required", errors.New("some fields are empty"))
		return
	}
	if len(req.Password) < 6 {
		respondWithError(w, http.StatusBadRequest, "Weak password", errors.New("password should be at least 6 charachters long"))
		return
	}
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unhasheable password", errors.New("bad password"))
		return
	}
	user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			respondWithError(w, http.StatusConflict, "Email already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create user account", err)
		return
	}
	accessToken, refreshToken, err := h.TokenService.GenerateTokens(r.Context(), user.ID.String(), "user_type", user.Email)
	if err != nil {
		respondWithError(w, 500, "failed to create tokens", err)
		return
	}
	tokenResoponse := TokenResoponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.TokenService.JWTConfig.AccessTokenExpiry.Seconds()),
	}
	respondWithJSON(w, http.StatusCreated, SignupResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		TokenResoponse: &tokenResoponse,
	})
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "bad request", err)
		return
	}
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "All fields are required", errors.New("some fields are empty"))
		return
	}
	user, err := h.DB.GetUserLogin(r.Context(), req.Email)

	if err != nil {
		respondWithError(w, 422, "Wrong email", err)
		return
	}

	valid, err := CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 422, "Wrong password", err)
		return
	}

	if !valid {
		respondWithError(w, 422, "Wrong password", errors.New("wrong password"))
		return
	}

	accessToken, refreshToken, err := h.TokenService.GenerateTokens(r.Context(), user.ID.String(), "user_type", user.Email)
	if err != nil {
		respondWithError(w, 500, "failed to create tokens", err)
		return
	}
	tokenResoponse := TokenResoponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.TokenService.JWTConfig.AccessTokenExpiry.Seconds()),
	}
	respondWithJSON(w, http.StatusOK, SignupResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		TokenResoponse: &tokenResoponse,
	})
}
