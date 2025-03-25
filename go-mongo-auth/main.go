package main

import (
	"log"
	"os"

	"example.com/go-mongo-auth/config"
	"example.com/go-mongo-auth/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Încarcă variabilele din fișierul .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Nu am putut încărca fișierul .env, se folosesc variabilele de mediu existente")
	}

	// Inițializează conexiunea cu MongoDB folosind variabila MONGODB_URI
	config.InitMongo()

	// Setează routerul Gin
	router := gin.Default()

	// Configurează rutele aplicației
	routes.SetupRoutes(router)

	// Obține portul din variabila de mediu, implicit 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serverul rulează la http://localhost:%s", port)
	router.Run(":" + port)
}