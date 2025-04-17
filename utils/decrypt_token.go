package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/heronhoga/memoraire-be/config"
	"github.com/heronhoga/memoraire-be/models"
	"github.com/joho/godotenv"
)

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		//decrypt token
		decryptedToken, err := Decrypt(tokenStr)
		if err != nil {
		    http.Error(w, "Unauthorized", http.StatusUnauthorized)
		    return
		}

		// Load secret
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
			return
		}
		jwtSecret := os.Getenv("JWT_KEY")

		// Parse and verify token
		token, err := jwt.Parse(decryptedToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			//check token in the database
			var existingUser models.User
			checkToken := config.DB.Where("username = ?", userID).Where("session = ?", tokenStr).First(&existingUser)

			if checkToken.RowsAffected == 0 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			issuer, ok := claims["iss"].(string)
			if !ok || issuer != "memoraire" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Save to context
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

