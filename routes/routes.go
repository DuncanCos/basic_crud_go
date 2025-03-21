package routes

import (
	"crud/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes configure toutes les routes de l'API
func SetupRoutes(router *gin.Engine, db *pgxpool.Pool) {
	// Création du contrôleur
	albumController := controllers.NewAlbumController(db)
	
	// Configuration des routes
	router.GET("/albums", albumController.GetAlbums)
	router.GET("/albums/:id", albumController.GetAlbumByID)
	router.POST("/albums", albumController.CreateAlbum)
	router.PUT("/albums/:id", albumController.UpdateAlbum)
	router.DELETE("/albums/:id", albumController.DeleteAlbum)
}