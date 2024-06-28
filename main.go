package main

import (
	"fmt"
	"natter-chat-go/db"
	"natter-chat-go/handlers"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:      "cquark",
		Passwd:    "cquark",
		DBName:    "natter",
		Addr:      "localhost:3306",
		ParseTime: true,
	})

	// check if there is an error connecting to the database
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// API routing
	mux := http.NewServeMux()

	handlers.ConfigurePostRoutes(mux, db)
	handlers.ConfigureAuthRoutes(mux, db)
	handlers.ConfigureLikesRoutes(mux, db)
	handlers.ConfigureCommentsRoutes(mux, db)

	corsMux := EnableCors(mux)

	if err := http.ListenAndServe("localhost:8080", corsMux); err != nil {
		fmt.Println(err.Error())
	}
}

// cors handler
func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
