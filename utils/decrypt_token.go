package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		//decrypt token
		decryptedJWT, err := Decrypt(token)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(decryptedJWT)

		next.ServeHTTP(w, r)
	})
}
