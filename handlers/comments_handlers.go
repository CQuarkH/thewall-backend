package handlers

import (
	"database/sql"
	"encoding/json"
	"natter-chat-go/models"
	"natter-chat-go/services"
	"net/http"
	"strconv"
)

func HandlePostComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+err.Error(), http.StatusBadRequest)
			return
		}

		comment := models.CreateCommentRequest{}
		err = json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			http.Error(w, "Error decoding comment: "+err.Error(), http.StatusBadRequest)
			return
		}

		comment.PostID = postID
		_, err = services.CreateComment(db, &comment)
		if err != nil {
			http.Error(w, "Error creating comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(comment)
	}
}

func HandleDeleteComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		commentID, err := strconv.Atoi(r.PathValue("commentID"))
		if err != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+err.Error(), http.StatusBadRequest)
			return
		}

		err = services.DeleteComment(db, commentID)
		if err != nil {
			http.Error(w, "Error deleting comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleGetPostComments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		postID, err := strconv.Atoi(r.PathValue("postID"))
		if err != nil {
			http.Error(w, "Invalid ID. Must be a positive number."+err.Error(), http.StatusBadRequest)
			return
		}

		comments, err := services.GetPostComments(db, postID)
		if err != nil {
			http.Error(w, "Error fetching comments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(comments)
	}
}

func ConfigureCommentsRoutes(router *http.ServeMux, db *sql.DB) {
	router.HandleFunc("POST /api/comments/{postID}", HandlePostComment(db))
	router.HandleFunc("DELETE /api/comments/{commentID}", HandleDeleteComment(db))
	router.HandleFunc("GET /api/comments/{postID}", HandleGetPostComments(db))
}
