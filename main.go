package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/vaxxnsh/metaverse/internals/database"
)

type apiConfig struct {
	DB *database.Queries
}

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetEnvVars(varName string) string {
	val := os.Getenv(varName)

	if val == "" {
		log.Fatalf("Can't find variable %s in Env", varName)
	}

	return val
}

func main() {
	godotenv.Load()
	port := GetEnvVars("PORT")
	dbUrl := GetEnvVars("DB_URL")

	ctx := context.Background()
	router := gin.Default()

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal("Error connecting with database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	fmt.Printf("The Server is running on the port: %s\n", port)
	router.Run(":" + port)
}
