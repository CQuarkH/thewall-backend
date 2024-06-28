package handlers

import (
	"database/sql"
	"encoding/json"
	"natter-chat-go/models"
	"natter-chat-go/services"
	"net/http"
	"strconv"
)

// get all posts
func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		posts, err := services.GetPosts(db)
		if err != nil {
			http.Error(w, "Error fetching all posts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(posts)
	}
}

func GetPostByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(r.PathValue("postID"))

		if err != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+err.Error(), http.StatusBadRequest)
			return
		}

		post, err := services.GetPostByID(db, id)
		if err != nil {
			http.Error(w, "Error fetching post by ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(post)
	}

}

func GetPostsByUsernameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		username := r.PathValue("username")

		user, err := services.GetUserByUsername(db, username)
		if err != nil {
			http.Error(w, "Error fetching user by username: "+err.Error(), http.StatusInternalServerError)
			return
		}

		posts, err := services.GetPostsByUserID(db, user.ID)
		if err != nil {
			http.Error(w, "Error fetching posts by user ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(posts)
	}

}

// create new post
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var post models.CreatePostRequest
		err := json.NewDecoder(r.Body).Decode(&post)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := ExtractUserFromToken(db, r)

		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusUnauthorized)
			return
		}

		post.UserID = user.ID

		createdPost, err := services.CreatePost(db, post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdPost)
	}
}

func UpdatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		postID, postErr := strconv.Atoi(r.PathValue("postID"))

		if postErr != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+postErr.Error(), http.StatusBadRequest)
			return
		}

		var post models.CreatePostRequest
		err := json.NewDecoder(r.Body).Decode(&post)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedPost, err := services.UpdatePost(db, post, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedPost)
	}
}

func DeletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(r.PathValue("postID"))

		if err != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+err.Error(), http.StatusBadRequest)
			return
		}

		_, err = services.DeletePost(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// configure post routes
func ConfigurePostRoutes(router *http.ServeMux, db *sql.DB) {
	router.HandleFunc("GET /api/posts", GetPostsHandler(db))
	router.HandleFunc("GET /api/post/{postID}", GetPostByIDHandler(db))
	router.HandleFunc("GET /api/posts/user/{username}", GetPostsByUsernameHandler(db))
	router.HandleFunc("POST /api/posts", CreatePostHandler(db))
	router.HandleFunc("PUT /api/posts/{postID}", UpdatePostHandler(db))
	router.HandleFunc("DELETE /api/posts/{postID}", DeletePostHandler(db))
}
