package services

import (
	"database/sql"
	"fmt"
	"natter-chat-go/models"
)

func ModifyPostLike(db *sql.DB, userID *int, postID *int, action string) error {
	var query string

	switch action {
	case "like":
		query = `
			INSERT INTO post_likes (user_id, post_id)
			VALUES (?, ?)
		`
	case "dislike":
		query = `
			DELETE FROM post_likes
			WHERE user_id = ? AND post_id = ?
		`
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	_, err := db.Exec(query, userID, postID)
	if err != nil {
		return err
	}

	return nil
}

func GetPostLikes(db *sql.DB, postID int) ([]int, error) {
	query := `
		SELECT user_id
		FROM post_likes
		WHERE post_id = ?
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	likes := []int{}
	for rows.Next() {
		var like int
		err := rows.Scan(&like)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

func GetLikedPostsByUserID(db *sql.DB, userID int) ([]int, error) {
	query := `
		SELECT post_id
		FROM post_likes
		WHERE user_id = ?
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	likedPosts := []int{}
	for rows.Next() {
		var postID int
		err := rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		likedPosts = append(likedPosts, postID)
	}

	return likedPosts, nil
}

func GetLikedPostsDetailsByUserID(db *sql.DB, userID int) ([]models.Post, error) {
	postIDs, err := GetLikedPostsByUserID(db, userID)
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	for _, postID := range postIDs {
		post, err := GetPostByID(db, postID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}

	return posts, nil
}
