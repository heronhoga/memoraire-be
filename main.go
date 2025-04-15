package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/heronhoga/memoraire-be/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	  return
	}
	config.DatabaseInit()

	fmt.Print("Memoraire backend is listening at port ", os.Getenv("APP_PORT"))
	http.ListenAndServe(":" + os.Getenv("APP_PORT"), nil)
}