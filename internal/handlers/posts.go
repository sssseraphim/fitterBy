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
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Content       string    `json:"content"`
	AuthorId      uuid.UUID `json:"author_id"`
	AuthorName    string    `json:"author_name"`
	MediaUrls     []string  `json:"media_urls"`
	Visibility    string    `json:"visibility"`
	LikesCount    int       `json:"likes_count"`
	Liked         bool      `json:"liked"`
	CommentsCount int       `json:"comments_count"`
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
	respondWithJSON(w, http.StatusCreated, Post{
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
	author, err := h.DB.GetUser(r.Context(), post.AuthorID.UUID)
	if err != nil {
		respondWithError(w, 404, "no author found", err)
		return
	}

	respondWithJSON(w, 200, Post{
		ID:            post.ID,
		CreatedAt:     post.CreatedAt,
		Content:       post.Content,
		MediaUrls:     post.MediaUrls,
		AuthorId:      post.AuthorID.UUID,
		AuthorName:    author.Name,
		Visibility:    post.Visibility,
		LikesCount:    int(post.LikeCount),
		CommentsCount: int(post.CommentsCount),
	})
}

func (h *PostHandler) HandleGetFollowedPosts(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)

	posts, err := h.DB.GetFollowedPosts(r.Context(), userId)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("failed to show posts: %v", err), err)
		return
	}
	var resp struct {
		Posts []Post `json:"posts"`
	}
	for _, p := range posts {
		fmt.Print(p)
		resp.Posts = append(resp.Posts, Post{
			ID:            p.ID,
			CreatedAt:     p.CreatedAt,
			Content:       p.Content,
			MediaUrls:     p.MediaUrls,
			AuthorId:      p.AuthorID,
			AuthorName:    p.AuthorName,
			Visibility:    p.Visibility,
			LikesCount:    int(p.LikeCount),
			CommentsCount: int(p.CommentsCount),
			Liked:         p.UserLiked,
		})
	}
	respondWithJSON(w, 200, resp)
}

func (h *PostHandler) HandlerLikePost(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	var req struct {
		PostId uuid.UUID `json:"post_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "wrong request", err)
		return
	}
	checkResult, err := h.DB.CheckUserLikedPost(r.Context(), database.CheckUserLikedPostParams{
		UserID: userId,
		PostID: req.PostId})
	if err != nil {
		respondWithError(w, 500, "failed to like post", err)
		return
	}
	var message string
	if !checkResult {
		err = h.DB.LikePost(r.Context(), database.LikePostParams{
			UserID: userId,
			PostID: req.PostId,
		})
		if err != nil {
			respondWithError(w, http.StatusConflict, "failed to like", err)
			return
		}
		message = "Post unliked"
	} else {
		err = h.DB.DeleteUserPostLike(r.Context(), database.DeleteUserPostLikeParams{
			UserID: userId,
			PostID: req.PostId,
		})
		if err != nil {
			respondWithError(w, http.StatusConflict, "failed to unlike", err)
			return
		}
		message = "Post unliked"
	}
	likeCount, _ := h.DB.GetLikeCount(r.Context(), req.PostId)
	respondWithJSON(w, http.StatusCreated, map[string]any{"success": true, "message": message, "likes_count": likeCount, "user_liked": !checkResult})
}

func userIdFromContext(r *http.Request) uuid.UUID {
	return uuid.MustParse(r.Context().Value(middleware.UserIDKey).(string))
}

func (h *PostHandler) HandlerComment(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	var req struct {
		PostId  uuid.UUID `json:"post_id"`
		Content string    `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "wrong request", err)
		return
	}
	err := h.DB.CommentOnPost(r.Context(), database.CommentOnPostParams{
		UserID:  userId,
		PostID:  req.PostId,
		Content: req.Content,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to comment", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"success": "success"})
}

func (h *PostHandler) HandlerGetComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		respondWithError(w, 400, "post_id query parameter is required", nil)
		return
	}
	postId, err := uuid.Parse(postIDStr)
	if err != nil {
		respondWithError(w, 400, "wrong post id format", err)
		return
	}
	comments, err := h.DB.GetPostComments(r.Context(), postId)
	if err != nil {
		respondWithError(w, 500, "failed to get comment", err)
		return
	}
	respondWithJSON(w, 200, comments)
}
