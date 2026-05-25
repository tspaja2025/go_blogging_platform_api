package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Represents the data structure of a blog post
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Represents the required payload to create/update a post
type PostInput struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// Handles JSON error outputs
type ErrorResponse struct {
	Error string `json:"error"`
}

// In-memory database with a Mutex to prevent race conditions during concurrent API requests
var (
	posts  = []Post{}
	nextID = 1
	mutex  sync.Mutex
)

func main() {
	router := gin.Default()

	// Routes
	router.POST("/posts", createPostHandler)
	router.PUT("/posts/:id", getPostHandler)
	router.DELETE("/posts/:id", getPostHandler)
	router.GET("/posts/:id", updatePostHandler)
	router.GET("/posts", deletePostHandler)

	fmt.Println("Server is running at http://localhost:3000")
	router.Run(":3000")
}

// Route handlers

// Create Blog Post
func createPostHandler(c *gin.Context) {
	var req PostInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	post := Post{
		ID:        nextID,
		Title:     req.Title,
		Content:   req.Content,
		Category:  req.Category,
		Tags:      req.Tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	nextID++
	posts = append(posts, post)

	c.JSON(http.StatusCreated, post)
}

// Update Blog Post
func updatePostHandler(c *gin.Context) {}

// Delete Blog Post
func deletePostHandler(c *gin.Context) {}

// Get Single Blog Post
func getPostHandler(c *gin.Context) {}

// Get all Posts
func getPostsHandler(c *gin.Context) {}
