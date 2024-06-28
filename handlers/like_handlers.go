package handlers

import (
	"database/sql"
	"encoding/json"
	"natter-chat-go/services"
	"net/http"
	"strconv"
)

// to get all posts liked by a user
func GetLikedPostsByUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID, err := strconv.Atoi(r.PathValue("userID"))
		if err != nil {
			http.Error(w, "Error parsing userID: "+err.Error(), http.StatusUnauthorized)
			return
		}

		posts, err := services.GetLikedPostsDetailsByUserID(db, userID)
		if err != nil {
			http.Error(w, "Error fetching liked posts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(posts)
	}
}

// to get all likes of a post
func GetPostLikesByPostIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		postID, err := strconv.Atoi(r.PathValue("postID"))

		if err != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+err.Error(), http.StatusBadRequest)
			return
		}

		likes, err := services.GetPostLikes(db, postID)
		if err != nil {
			http.Error(w, "Error fetching all likes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(likes)
	}
}

// to like or dislike a post
func ModifyLikeHandler(db *sql.DB, action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		postID, err := strconv.Atoi(r.PathValue("postID"))

		if err != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+err.Error(), http.StatusBadRequest)
			return
		}

		user, err := ExtractUserFromToken(db, r)
		if err != nil {
			http.Error(w, "Error extracting user from token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = services.ModifyPostLike(db, &user.ID, &postID, action)
		if err != nil {
			http.Error(w, "Error modifying post like: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("Post like modified successfully")
	}
}

func ConfigureLikesRoutes(router *http.ServeMux, db *sql.DB) {
	router.HandleFunc("GET /api/likes/by-user/{userID}", GetLikedPostsByUserHandler(db))
	router.HandleFunc("GET /api/likes/{postID}", GetPostLikesByPostIDHandler(db))
	router.HandleFunc("PUT /api/likes/{postID}/like", ModifyLikeHandler(db, "like"))
	router.HandleFunc("PUT /api/likes/{postID}/dislike", ModifyLikeHandler(db, "dislike"))
}
