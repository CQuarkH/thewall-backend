package handlers

import (
	"database/sql"
	"encoding/json"
	"natter-chat-go/models"
	"natter-chat-go/services"
	"net/http"
	"time"
)

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := services.Login(db, creds.Email, creds.Password)
		if err != nil {
			http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
			return
		}

		// set the token in a cookie
		cookie := http.Cookie{
			Name:     "token",
			Value:    user.Token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		}

		user.Token = ""

		http.SetCookie(w, &cookie)
		json.NewEncoder(w).Encode(user)
	}
}

func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now().AddDate(0, 0, -1),
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		}

		http.SetCookie(w, &cookie)
	}
}

// user register
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var register models.Register
		err := json.NewDecoder(r.Body).Decode(&register)
		if err != nil {
			http.Error(w, "Error decoding json request: "+err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := services.Register(db, register)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(userID)
	}
}

func ConfigureAuthRoutes(muxRouter *http.ServeMux, db *sql.DB) {
	muxRouter.HandleFunc("POST /api/auth/login", Login(db))
	muxRouter.HandleFunc("POST /api/auth/register", Register(db))
	muxRouter.HandleFunc("POST /api/auth/logout", Logout(db))
}
