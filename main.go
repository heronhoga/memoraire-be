package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/heronhoga/memoraire-be/config"
	"github.com/heronhoga/memoraire-be/routes"
	"github.com/joho/godotenv"
)

func main() {
	httpServer := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	  return
	}
	config.DatabaseInit()
	
	//routes
	routes.UserRoutes(httpServer)

	fmt.Println("Memoraire backend is listening at port: ", os.Getenv("APP_PORT"))
	http.ListenAndServe(":" + os.Getenv("APP_PORT"), httpServer)
}