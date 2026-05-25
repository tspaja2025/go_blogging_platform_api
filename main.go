package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
func updatePostHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid post ID",
		})
		return
	}

	var req PostInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, post := range posts {
		if post.ID == id {
			posts[i].Title = req.Title
			posts[i].Content = req.Content
			posts[i].Category = req.Category
			posts[i].Tags = req.Tags
			posts[i].UpdatedAt = time.Now()

			c.JSON(http.StatusOK, posts[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "post not found",
	})
}

// Delete Blog Post
func deletePostHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid post ID",
		})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, post := range posts {
		if post.ID == id {
			posts = append(posts[:i], posts[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "post not found",
	})
}

// Get Single Blog Post
func getPostHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid post ID",
		})
		return
	}

	for _, post := range posts {
		if post.ID == id {
			c.JSON(http.StatusOK, post)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "post not found",
	})
}

// Get all Posts
func getPostsHandler(c *gin.Context) {
	term := strings.ToLower(c.Query("term"))

	if term == "" {
		c.JSON(http.StatusOK, posts)
		return
	}

	filteredPosts := []Post{}

	for _, post := range posts {
		if strings.Contains(strings.ToLower(post.Title), term) ||
			strings.Contains(strings.ToLower(post.Content), term) ||
			strings.Contains(strings.ToLower(post.Category), term) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	c.JSON(http.StatusOK, filteredPosts)
}
