package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sssseraphim/fitterBy/internal/database"
	"github.com/sssseraphim/fitterBy/internal/middleware"
)

type UserHandler struct {
	DB *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	Premium   bool      `json:"premium"`
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle by id")
	userIdString := r.PathValue("user_id")
	fmt.Println(userIdString)
	if userIdString == "" {
		respondWithError(w, 400, "user id required", errors.New("no id"))
		return
	}
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		respondWithError(w, 400, "incorrect id", err)
		return
	}

	user, err := h.DB.GetUser(r.Context(), userId)
	if err != nil {
		respondWithError(w, 404, "no user found", err)
		return
	}
	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Bio:       user.Bio,
		Premium:   user.Premium.Bool,
	})
}

func (h *UserHandler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle me")

	userIdString := r.Context().Value(middleware.UserIDKey).(string)
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		respondWithError(w, 400, "incorrect id", err)
		return
	}

	user, err := h.DB.GetUser(r.Context(), userId)
	if err != nil {
		respondWithError(w, 404, "no user found", err)
		return
	}
	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Name:      user.Name,
		Bio:       user.Bio,
		Premium:   user.Premium.Bool,
	})
}

func (h *UserHandler) HandlerUpdateBio(w http.ResponseWriter, r *http.Request) {
	fmt.Println("adding bio")

	userIdString := r.Context().Value(middleware.UserIDKey).(string)
	userId := uuid.MustParse(userIdString)

	var request struct {
		Bio string `json:"bio"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, 400, "failed to decoder the bio", err)
		return
	}
	err := h.DB.UpdateBio(r.Context(), database.UpdateBioParams{
		ID:  userId,
		Bio: request.Bio})
	if err != nil {
		respondWithError(w, 500, "failed to update the bio", err)
		return
	}
	respondWithJSON(w, 200, request)
}
