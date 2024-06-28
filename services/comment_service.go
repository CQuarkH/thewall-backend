package services

import (
	"database/sql"
	"natter-chat-go/models"
	"time"
)

func CreateComment(db *sql.DB, comment *models.CreateCommentRequest) (*models.Comment, error) {
	comment.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	query := `
		INSERT INTO comments (content, created_at, user_id, post_id)
		VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(query, comment.Content, comment.CreatedAt, comment.UserID, comment.PostID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Comment{
		ID:        int(id),
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
	}, nil
}

func DeleteComment(db *sql.DB, commentID int) error {
	query := `
		DELETE FROM comments
		WHERE id = ?
	`

	_, err := db.Exec(query, commentID)
	return err

}

func GetPostComments(db *sql.DB, postID int) ([]*models.CommentWithUserResponse, error) {
	query := `
		SELECT id, content, created_at, user_id, post_id
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at DESC
	`

	return getComments(db, query, postID)
}

func GetUserComments(db *sql.DB, userID int) ([]*models.CommentWithUserResponse, error) {
	query := `
		SELECT id, content, created_at, user_id, post_id
		FROM comments
		WHERE user_id = ?
	`
	return getComments(db, query, userID)
}

// generic function to get comments
func getComments(db *sql.DB, query string, args ...interface{}) ([]*models.CommentWithUserResponse, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.CommentWithUserResponse{}
	for rows.Next() {
		var comment models.CommentWithUserResponse
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.UserID, &comment.PostID); err != nil {
			return nil, err
		}

		user, err := GetUserByID(db, comment.UserID)
		if err != nil {
			return nil, err
		}

		comment.Username = user.Username
		comment.UserPhoto = user.Photo

		comments = append(comments, &comment)
	}

	return comments, nil
}
