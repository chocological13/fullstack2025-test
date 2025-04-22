package client

import (
	"encoding/json"
	"fmt"
	db "fullstack2025-test/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
)

type ClientHandler struct {
	db  *db.Queries
	rdb *redis.Client
}

func NewClientHandler(db *db.Queries, rdb *redis.Client) *ClientHandler {
	return &ClientHandler{db: db, rdb: rdb}
}

func (h *ClientHandler) ListClients(c *gin.Context) {
	clients, err := h.db.ListClients(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"clients": clients})
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	var params db.CreateClientParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.db.CreateClient(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store in Redis
	clientJSON, _ := json.Marshal(client)
	h.rdb.Set(c, fmt.Sprintf("client:%s", client.Slug), clientJSON, 0)

	c.JSON(http.StatusCreated, gin.H{"client": client})
}

func (h *ClientHandler) GetClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	client, err := h.db.GetClient(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client": client})
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {
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

	client, err := h.db.UpdateClient(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update Redis
	clientJSON, _ := json.Marshal(client)
	h.rdb.Set(c, fmt.Sprintf("client:%s", client.Slug), clientJSON, 0)

	c.JSON(http.StatusOK, gin.H{"client": client})
}

func (h *ClientHandler) DeleteClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	// Get client slug before deletion
	client, err := h.db.GetClient(c, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	err = h.db.DeleteClient(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete from Redis
	h.rdb.Del(c, fmt.Sprintf("client:%s", client.Slug))

	c.JSON(http.StatusOK, gin.H{"delete client": "success"})
}
