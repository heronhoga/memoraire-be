package utils

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	appKey string
	once   sync.Once
)

func loadEnv() {
	_ = godotenv.Load()
	appKey = os.Getenv("APP_KEY")
	if appKey == "" {
		log.Println("Warning: APP_KEY not set")
	}
}

func CheckKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		once.Do(loadEnv)

		keyHeader := r.Header.Get("hgtoken")
		log.Println("token from frontend: ", keyHeader)
		if keyHeader == "" || keyHeader != appKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
