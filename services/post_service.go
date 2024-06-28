package services

import (
	"database/sql"
	"log"
	"natter-chat-go/models"
	"strings"
	"time"
)

func scanPost(rows *sql.Rows) (models.Post, error) {
	var post models.Post
	var photoURLs string

	err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UserID, &photoURLs, &post.CommentCount)
	if err != nil {
		return post, err
	}

	if photoURLs != "" {
		post.PhotoURLs = strings.Split(photoURLs, ",")
	} else {
		post.PhotoURLs = []string{}
	}

	return post, nil
}

func populatePostDetails(db *sql.DB, post *models.Post) error {
	var err error
	post.LikedBy, err = GetPostLikes(db, post.ID)
	if err != nil {
		return err
	}

	userProfile, err := GetUserProfileByID(db, post.UserID)
	if err != nil {
		return err
	}

	post.User = *userProfile

	return nil
}

func GetPosts(db *sql.DB) ([]models.Post, error) {
	var posts []models.Post

	query := `
		SELECT p.id, p.title, p.content, p.created_at, p.user_id, 
		       COALESCE(GROUP_CONCAT(ph.url), '') AS photo_urls,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count
		FROM posts p
		LEFT JOIN photos ph ON p.id = ph.post_id
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}

		err = populatePostDetails(db, &post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostByID(db *sql.DB, id int) (*models.Post, error) {
	var post models.Post
	var photoURLs string

	query := `
		WITH photo_urls AS (
			SELECT p.id, COALESCE(GROUP_CONCAT(ph.url), '') AS photo_urls
			FROM posts p
			LEFT JOIN photos ph ON p.id = ph.post_id
			WHERE p.id = ?
			GROUP BY p.id
		),
		comment_counts AS (
			SELECT p.id, COUNT(c.id) AS comment_count
			FROM posts p
			LEFT JOIN comments c ON p.id = c.post_id
			WHERE p.id = ?
			GROUP BY p.id
		)
		SELECT p.id, p.title, p.content, p.created_at, p.user_id, 
		       pu.photo_urls, cc.comment_count
		FROM posts p
		LEFT JOIN photo_urls pu ON p.id = pu.id
		LEFT JOIN comment_counts cc ON p.id = cc.id
		WHERE p.id = ?
	`

	err := db.QueryRow(query, id, id, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UserID, &photoURLs, &post.CommentCount)
	if err != nil {
		return nil, err
	}

	if photoURLs != "" {
		post.PhotoURLs = strings.Split(photoURLs, ",")
	} else {
		post.PhotoURLs = []string{}
	}

	err = populatePostDetails(db, &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func GetPostsByUserID(db *sql.DB, userID int) ([]models.Post, error) {
	var posts []models.Post

	query := `
		WITH photo_urls AS (
			SELECT p.id, COALESCE(GROUP_CONCAT(ph.url), '') AS photo_urls
			FROM posts p
			LEFT JOIN photos ph ON p.id = ph.post_id
			GROUP BY p.id
		),
		comment_counts AS (
			SELECT p.id, COUNT(c.id) AS comment_count
			FROM posts p
			LEFT JOIN comments c ON p.id = c.post_id
			GROUP BY p.id
		)
		SELECT p.id, p.title, p.content, p.created_at, p.user_id, 
		       pu.photo_urls, cc.comment_count
		FROM posts p
		LEFT JOIN photo_urls pu ON p.id = pu.id
		LEFT JOIN comment_counts cc ON p.id = cc.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var photoURLs string

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UserID, &photoURLs, &post.CommentCount)
		if err != nil {
			return nil, err
		}

		if photoURLs != "" {
			post.PhotoURLs = strings.Split(photoURLs, ",")
		} else {
			post.PhotoURLs = []string{}
		}

		err = populatePostDetails(db, &post)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(db *sql.DB, post models.CreatePostRequest) (*models.Post, error) {
	result, err := db.Exec("INSERT INTO posts (title, content, created_at, user_id) VALUES (?, ?, ?, ?)", post.Title, post.Content, time.Now(), post.UserID)
	if err != nil {
		return &models.Post{}, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return &models.Post{}, err
	}

	for _, photoURL := range post.PhotoURLs {
		_, err := db.Exec("INSERT INTO photos (url, post_id) VALUES (?, ?)", photoURL, lastInsertId)
		if err != nil {
			return &models.Post{}, err
		}
	}

	return GetPostByID(db, int(lastInsertId))
}

func UpdatePost(db *sql.DB, post models.CreatePostRequest, postID int) (*models.Post, error) {
	// update photos (delete all photos and insert the new ones)
	if len(post.PhotoURLs) > 0 {
		_, err := db.Exec("DELETE FROM photos WHERE post_id = ?", postID)
		if err != nil {
			return &models.Post{}, err
		}

		for _, photoURL := range post.PhotoURLs {
			_, err := db.Exec("INSERT INTO photos (url, post_id) VALUES (?, ?)", photoURL, postID)
			if err != nil {
				return &models.Post{}, err
			}
		}

	}
	_, err := db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", post.Title, post.Content, postID)
	if err != nil {
		return &models.Post{}, err
	}

	return GetPostByID(db, postID)
}

func DeletePost(db *sql.DB, id int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("rollback due to panic: %v", r)
		}
	}()

	// delete photos
	if _, err := tx.Exec("DELETE FROM photos WHERE post_id = ?", id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// delete likes
	if _, err := tx.Exec("DELETE FROM post_likes WHERE post_id = ?", id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// delete comments
	if _, err := tx.Exec("DELETE FROM comments WHERE post_id = ?", id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// delete post
	result, err := tx.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
