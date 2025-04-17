package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)
func GenerateJWT(userId string) (string, error) {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	  return "", err
	}

	//jwt secret
	var jwtSecret = os.Getenv("JWT_KEY")

	//set claims
	claims := jwt.MapClaims{
		"user_id": userId,
		"iat":     time.Now().Unix(),              
		"iss":     "memoraire",                       
	}

	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, errSignToken := token.SignedString([]byte(jwtSecret))

	if errSignToken != nil {
		return "", errSignToken
	}

	return signedToken, nil
}