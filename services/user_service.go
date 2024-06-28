package services

import (
	"database/sql"
	"natter-chat-go/models"
)

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, email, photo_url FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Email, &user.Photo)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, email, photo_url FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Photo)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserProfileByID(db *sql.DB, id int) (*models.UserProfile, error) {
	var user models.UserProfile

	query := "SELECT username, COALESCE(photo_url, '') FROM users WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&user.Username, &user.Photo)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
