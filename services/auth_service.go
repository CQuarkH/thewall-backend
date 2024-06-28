package services

import (
	"database/sql"
	"errors"
	"natter-chat-go/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, email, password string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, email, password, photo_url FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo)
	if err != nil {
		return nil, err
	}

	// password verification
	err = verifyPassword(password, user.Password)
	if err != nil {
		return nil, err
	}

	// no return password
	user.Password = ""
	token, err := CreateJWT(email)
	if err != nil {
		return nil, err
	}

	user.Token = token

	return &user, nil
}

func Register(db *sql.DB, register models.Register) (int64, error) {

	// check if email already exists
	var email string
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", register.Email).Scan(&email)
	if err == nil {
		return 0, errors.New("email already exists")
	}

	hashedPassword, err := hashPassword(register.Password)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec("INSERT INTO users (username, email, password, photo_url) VALUES (?, ?, ?, ?)", register.Username, register.Email, hashedPassword, register.Photo)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// jwt features
func CreateJWT(email string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = expiration.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("natter-secret-key"))
}
