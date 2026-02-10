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

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/users", func(c *gin.Context) {
		users, err := apiCfg.DB.ListUsers(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"count":   len(users),
			"users":   users,
		})
	})

	router.POST("/users", func(c *gin.Context) {
		var req createUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		now := time.Now()
		user, err := apiCfg.DB.CreateUser(c.Request.Context(), database.CreateUserParams{
			ID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			Name:     req.Name,
			Password: req.Password,
			CreatedAt: pgtype.Timestamp{
				Time:  now,
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  now,
				Valid: true,
			},
		})

		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"user":    user,
		})
	})

	fmt.Printf("The Server is running on the port: %s\n", port)
	router.Run(":" + port)
}
