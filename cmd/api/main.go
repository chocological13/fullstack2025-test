package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "fullstack2025-test/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"

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

	// Routes
	apiGroup := router.Group("/api")
	{
		clients := apiGroup.Group("/clients")
		{
			clients.GET("/", app.listClients)
			clients.POST("/", app.createClient)
			clients.GET("/:id", app.getClient)
			clients.PUT("/:id", app.updateClient)
			clients.DELETE("/:id", app.deleteClient)
		}
	}

	// Start server
	fmt.Println("Server started on :8080", app)
	log.Fatal(router.Run(":8080"))

}

// Handlers
func (app *application) listClients(c *gin.Context) {
	clients, err := app.db.ListClients(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"clients": clients})
}

func (app *application) createClient(c *gin.Context) {
	var params db.CreateClientParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := app.db.CreateClient(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store in Redis
	clientJSON, _ := json.Marshal(client)
	app.redis.Set(c, fmt.Sprintf("client:%s", client.Slug), clientJSON, 0)

	c.JSON(http.StatusCreated, gin.H{"client": client})
}

func (app *application) getClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	client, err := app.db.GetClient(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client": client})
}

func (app *application) updateClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	var params db.UpdateClientParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.ID = int32(id)

	client, err := app.db.UpdateClient(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update Redis
	clientJSON, _ := json.Marshal(client)
	app.redis.Set(c, fmt.Sprintf("client:%s", client.Slug), clientJSON, 0)

	c.JSON(http.StatusOK, gin.H{"client": client})
}

func (app *application) deleteClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	// Get client slug before deletion
	client, err := app.db.GetClient(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	err = app.db.DeleteClient(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete from Redis
	app.redis.Del(c, fmt.Sprintf("client:%s", client.Slug))

	c.JSON(http.StatusOK, gin.H{"delete client": "success"})
}
