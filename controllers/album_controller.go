package controllers

import (
	"context"
	"net/http"

	"crud/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AlbumController gère les opérations CRUD pour les albums
type AlbumController struct {
	DB *pgxpool.Pool
}

// NewAlbumController crée une nouvelle instance de AlbumController
func NewAlbumController(db *pgxpool.Pool) *AlbumController {
	return &AlbumController{DB: db}
}

// GetAlbums récupère tous les albums
func (c *AlbumController) GetAlbums(ctx *gin.Context) {
	var albums []models.Album
	
	rows, err := c.DB.Query(context.Background(), "SELECT * FROM albums")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var album models.Album
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		albums = append(albums, album)
	}
	
	ctx.IndentedJSON(http.StatusOK, albums)
}

// GetAlbumByID récupère un album par son ID
func (c *AlbumController) GetAlbumByID(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var album models.Album
	err := c.DB.QueryRow(context.Background(), "SELECT * FROM albums WHERE id = $1", id).
		Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.IndentedJSON(http.StatusOK, gin.H{"result": album})
}

// CreateAlbum crée un nouvel album
func (c *AlbumController) CreateAlbum(ctx *gin.Context) {
	var newAlbum models.NewAlbum
	
	if err := ctx.BindJSON(&newAlbum); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	_, err := c.DB.Exec(context.Background(), 
		"INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3)", 
		newAlbum.Title, newAlbum.Artist, newAlbum.Price)
		
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.IndentedJSON(http.StatusCreated, newAlbum)
}

// UpdateAlbum met à jour un album existant
func (c *AlbumController) UpdateAlbum(ctx *gin.Context) {
	id := ctx.Param("id")
	var newAlbum models.NewAlbum
	
	if err := ctx.BindJSON(&newAlbum); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	_, err := c.DB.Exec(context.Background(), 
		"UPDATE albums SET title=$1, artist=$2, price=$3 WHERE id=$4", 
		newAlbum.Title, newAlbum.Artist, newAlbum.Price, id)
		
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.IndentedJSON(http.StatusCreated, newAlbum)
}

// DeleteAlbum supprime un album par son ID
func (c *AlbumController) DeleteAlbum(ctx *gin.Context) {
	id := ctx.Param("id")
	
	_, err := c.DB.Exec(context.Background(), "DELETE FROM albums WHERE id=$1", id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Album supprimé avec succès"})
}