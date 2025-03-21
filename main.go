package main

import (
	"fmt"
	"log"

	"crud/db"
	"crud/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Chargement des variables d'environnement
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialisation de la base de données
	dbPool, err := db.InitDB()
	if err != nil {
		fmt.Printf("Erreur lors de la connexion à la base de données: %v\n", err)
		return
	}
	defer dbPool.Close()

	// Initialisation du router
	router := gin.Default()
	
	// Configuration des routes
	routes.SetupRoutes(router, dbPool)

	// Démarrage du serveur
	router.Run("localhost:8080")
}