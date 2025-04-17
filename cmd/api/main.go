package main

import (
	"database/sql"
	"fmt"
	db "fullstack2025-test/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	db    *db.Queries
	redis *redis.Client
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	dsn := os.Getenv("DATABASE_URL")

	log.Println("dsn", dsn)

	// Connect to DB
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := db.New(dbConn)
	log.Println("Connected to database")

	// Connect to Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	opt, _ := redis.ParseURL(redisAddr)
	client := redis.NewClient(opt)
	log.Println("Connected to redis")

	// Initialize app
	app := &application{
		db:    dbQueries,
		redis: client,
	}

	router := gin.Default()

	// Start server
	fmt.Println("Server started on :8080", app)
	log.Fatal(router.Run(":8080"))
}
