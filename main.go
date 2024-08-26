package main

import (
	"NotesService/auntification"
	"NotesService/notes"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"strconv"
)

var (
	session     *auntification.Session
	redisClient *redis.Client
)

func main() {
	db, err := notes.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authService := auntification.NewAuthService(db)

	redisClient, err = auntification.SetupRedis("localhost:6379", "", 0)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	session = auntification.NewSession(redisClient)

	http.HandleFunc("/auth", authService.AuthenticateHandler)
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		sessionID := cookie.Value
		user, err := session.GetUser(sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Использовать информацию о пользователе
		r.Header.Set("UserID", strconv.Itoa(user.ID))
	})

	log.Println("Server is running on :9080")
	log.Fatal(http.ListenAndServe(":9080", nil))
}
