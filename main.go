package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/heronhoga/memoraire-be/config"
	"github.com/heronhoga/memoraire-be/routes"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	httpServer := http.NewServeMux()
	_ = godotenv.Load()
	config.DatabaseInit()
	
	//routes
	routes.UserRoutes(httpServer)
	routes.MemoRoutes(httpServer)

	//cors
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://memoraire.web.id", "https://localhost:3000", "http://memoraire.web.id", "http://memoraire.web.id:3030", "https://memoraire.web.id:3030"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "hgtoken"},
		AllowCredentials: true,
	})

	httpWithCors := corsHandler.Handler(httpServer)

	fmt.Println("Memoraire backend is listening at port:", os.Getenv("APP_PORT"))
	http.ListenAndServe(":" + os.Getenv("APP_PORT"), httpWithCors)
}