package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB initialise la connexion à la base de données
func InitDB() (*pgxpool.Pool, error) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("PASSWORD")
	db_name := os.Getenv("DBNAME")
	
	dbURL := "postgres://" + username + ":" + password + "@localhost:5432/" + db_name
	
	return pgxpool.New(context.Background(), dbURL)
}