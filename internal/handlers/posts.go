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

type PostHandler struct {
	DB *database.Queries
}

type Post struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Content    string    `json:"content"`
	AuthorId   uuid.UUID `json:"author_id"`
	AuthorName string    `json:"author_name"`
	MediaUrls  []string  `json:"media_urls"`
	Visibility string    `json:"visibility"`
}

func (h *PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Content    string   `json:"content"`
		MediaUrls  []string `json:"media_urls"`
		Visibility string   `json:"visibility"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong request", err)
		return
	}
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))
	post, err := h.DB.CreatePost(r.Context(), database.CreatePostParams{
		UserID:     userId,
		Content:    request.Content,
		MediaUrls:  request.MediaUrls,
		Visibility: request.Visibility,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create post", err)
		return
	}
	respondWithJSON(w, 200, Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Content:   post.Content,
		MediaUrls: post.MediaUrls,
		AuthorId:  post.UserID,
	})

}
func (h *PostHandler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle post by id")
	postIdString := r.PathValue("post_id")
	fmt.Println(postIdString)
	if postIdString == "" {
		respondWithError(w, 400, "post id required", errors.New("no id"))
		return
	}
	postId, err := uuid.Parse(postIdString)
	if err != nil {
		respondWithError(w, 400, "incorrect id", err)
		return
	}

	post, err := h.DB.GetPost(r.Context(), postId)
	if err != nil {
		respondWithError(w, 404, "no post found", err)
		return
	}
	if post.Visibility != "public" {
		respondWithError(w, http.StatusUnauthorized, "post is for followers only", errors.New("wrong post visibility"))
		return
	}
	author, err := h.DB.GetUser(r.Context(), post.UserID)
	if err != nil {
		respondWithError(w, 404, "no author found", err)
		return
	}

	respondWithJSON(w, 200, Post{
		ID:         post.ID,
		CreatedAt:  post.CreatedAt,
		Content:    post.Content,
		MediaUrls:  post.MediaUrls,
		AuthorId:   post.UserID,
		AuthorName: author.Name,
	})
}

func (h *PostHandler) HandleGetFollowedPosts(w http.ResponseWriter, r *http.Request) {
	userId := uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))

	posts, err := h.DB.GetFollowedPosts(r.Context(), userId)
	if err != nil {
		respondWithError(w, 404, "no posts found", err)
		return
	}
	var resp struct {
		Posts []Post `json:"posts"`
	}
	for _, p := range posts {
		resp.Posts = append(resp.Posts, Post{
			ID:         p.ID,
			CreatedAt:  p.CreatedAt,
			Content:    p.Content,
			MediaUrls:  p.MediaUrls,
			AuthorId:   p.AuthorID,
			AuthorName: p.AuthorName,
		})
	}
	respondWithJSON(w, 200, resp)
}
