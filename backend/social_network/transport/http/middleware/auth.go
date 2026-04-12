package middleware

import (
	"social_network/internal/auth"
	"social_network/internal/post"
	"social_network/transport/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*auth.Claims)
			if !ok || claims == nil {
				response.Error(w, ErrUnauthorized)
				return
			}

			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Error(w, ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequirePostOwner(postService post.PostService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserID(r.Context())
			if !ok {
				response.Error(w, ErrUnauthorized)
			}
			postID, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				response.Error(w, ErrInvalidID)
			}

			post, err := postService.GetByID(r.Context(), postID)
			if err != nil || post.UserID != userID {
				response.Error(w, ErrNotOwner)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
