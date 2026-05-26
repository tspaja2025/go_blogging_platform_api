package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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
	mutex  sync.RWMutex
)

func main() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("POST /posts", createPostHandler)
	mux.HandleFunc("GET /posts", getPostsHandler)
	mux.HandleFunc("GET /posts/{id}", getPostHandler)
	mux.HandleFunc("PUT /posts/{id}", updatePostHandler)
	mux.HandleFunc("DELETE /posts/{id}", deletePostHandler)

	fmt.Println("Server is running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}

// Helper functions
func validateInput(input PostInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("title is required and cannot be empty")
	}
	if strings.TrimSpace(input.Content) == "" {
		return fmt.Errorf("content is required and cannot be empty")
	}
	if strings.TrimSpace(input.Category) == "" {
		return fmt.Errorf("category is required and cannot be empty")
	}
	return nil
}

// Shortcut helper to write JSON responses safely
func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Formats and sends a bad request or error payload
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, ErrorResponse{Error: message})
}

// Route handlers

// Create Blog Post
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input PostInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := validateInput(input); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()
	post := Post{
		ID:        nextID,
		Title:     input.Title,
		Content:   input.Content,
		Category:  input.Category,
		Tags:      input.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	nextID++
	posts = append(posts, post)

	sendJSON(w, http.StatusCreated, post)
}

// Update Blog Post
func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	var input PostInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := validateInput(input); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, post := range posts {
		if post.ID == id {
			posts[i].Title = input.Title
			posts[i].Content = input.Content
			posts[i].Category = input.Category
			posts[i].Tags = input.Tags
			posts[i].UpdatedAt = time.Now()

			sendJSON(w, http.StatusOK, posts[i])
			return
		}
	}

	sendError(w, http.StatusNotFound, "Blog post not found")
}

// Delete Blog Post
func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, post := range posts {
		if post.ID == id {
			// Remove post from slice
			posts = append(posts[:i], posts[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	sendError(w, http.StatusNotFound, "Blog post not found")
}

// Get Single Blog Post
func getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	mutex.RLock()
	defer mutex.RUnlock()

	for _, post := range posts {
		if post.ID == id {
			sendJSON(w, http.StatusOK, post)
			return
		}
	}

	sendError(w, http.StatusNotFound, "Blog post not found")
}

// Get all Posts
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	term := strings.ToLower(r.URL.Query().Get("term"))

	mutex.RLock()
	defer mutex.RUnlock()

	// If no filter term is present, return all posts
	if term == "" {
		sendJSON(w, http.StatusOK, posts)
		return
	}

	filteredPosts := []Post{}
	for _, post := range posts {
		if contains(post.Title, term) ||
			contains(post.Content, term) ||
			contains(post.Category, term) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	sendJSON(w, http.StatusOK, filteredPosts)
}

func contains(text, term string) bool {
	return strings.Contains(strings.ToLower(text), term)
}
