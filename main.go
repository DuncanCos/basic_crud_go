package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// album struct .
type NewAlbum struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type Albums struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type App struct {
	DB *pgxpool.Pool
}

// main qui creer les routes
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("PASSWORD")
	db_name := os.Getenv("DBNAME")

	urlExample := "postgres://" + username + ":" + password + "@localhost:5432/" + db_name

	conn, err := pgxpool.New(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	app := &App{DB: conn}
	defer conn.Close()

	router := gin.Default()
	router.GET("/albums", app.getAlbums)
	router.GET("/albums/:id", app.getAlbumByID)
	router.POST("/albums", app.postAlbums)
	router.PUT("/albums/:id", app.putAlbums)
	router.DELETE("/albums/:id", app.deleteAlbums)

	router.Run("localhost:8080")
}

func (app *App) deleteAlbums(c *gin.Context) {
	// supprimer un element dans une table
	id := c.Param("id")
	_, err := app.DB.Exec(context.Background(), "DELETE FROM albums WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album supprimé avec succès"})
}

// fait tout les albums en json
func (app *App) getAlbums(c *gin.Context) {

	// retourne tout dune table
	var albums []Albums
	rows, err := app.DB.Query(context.Background(), "SELECT * FROM albums")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()
	for rows.Next() {
		var album Albums
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		albums = append(albums, album)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func (app *App) putAlbums(c *gin.Context) {
	id := c.Param("id")
	var newAlbum NewAlbum

	//c.bindjson qui creer newalbum et ce base sur la struct
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	fmt.Print(newAlbum)

	// update un element dans une table put
	_, err := app.DB.Exec(context.Background(), "UPDATE albums SET title=$1, artist=$2, price=$3 WHERE id=$4", newAlbum.Title, newAlbum.Artist, newAlbum.Price, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// va creer un album recu avec le json dans le body (ce base sur la struct)
func (app *App) postAlbums(c *gin.Context) {
	var newAlbum NewAlbum

	//c.bindjson qui creer newalbum et ce base sur la struct
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// ajoute un element dans une table post
	_, err := app.DB.Exec(context.Background(), "INSERT INTO albums (title, artist, price) VALUES ($1,$2,$3)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// fait un get qui ce base sur l'id du  parametre
func (app *App) getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// va retourner un seul object
	var album Albums
	err := app.DB.QueryRow(context.Background(), "select * from albums where id = $1", id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"result": album})
}
