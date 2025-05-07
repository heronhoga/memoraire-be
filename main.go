package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/heronhoga/memoraire-be/config"
	"github.com/heronhoga/memoraire-be/routes"
	"github.com/joho/godotenv"
)

func main() {
	httpServer := http.NewServeMux()
	_ = godotenv.Load()
	config.DatabaseInit()
	
	//routes
	routes.UserRoutes(httpServer)
	routes.MemoRoutes(httpServer)

	fmt.Println("Memoraire backend is listening at port:", os.Getenv("APP_PORT"))
	http.ListenAndServe(":" + os.Getenv("APP_PORT"), httpServer)
}