package handlers

import (
	"database/sql"
	"errors"
	"natter-chat-go/models"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func GetTokenFromCookies(r *http.Request) (string, error) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			return cookie.Value, nil
		}
	}

	return "", errors.New("cookie token not found")
}

func ExtractUserFromToken(db *sql.DB, r *http.Request) (models.User, error) {
	tokenString, err := GetTokenFromCookies(r)

	if err != nil {
		return models.User{}, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("natter-secret-key"), nil
	})

	if err != nil {
		return models.User{}, err
	}

	if !token.Valid {
		return models.User{}, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.User{}, err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return models.User{}, err
	}

	var user models.User
	err = db.QueryRow("SELECT id, username, email FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, err
		}
		return models.User{}, err
	}

	return user, nil
}
