package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"golang-multi-tenant/internal/database"
	"golang-multi-tenant/internal/models"
)

// @Summary     Create a new post
// @Description Create a new post for the authenticated user
// @Tags        posts
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body models.CreatePostRequest true "Post details"
// @Success     201 {object} models.Post "Post created successfully"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /posts [post]
func CreatePost(c *gin.Context) {
    var req models.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get user info from context (set by auth middleware)
    userID := c.GetInt("user_id")
    tenantID := c.GetInt("tenant_id")

    // Get tenant database
    tenantDB, err := database.GetTenantDB(tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // Create post
    var post models.Post
    err = tenantDB.QueryRow(`
        INSERT INTO posts (user_id, title, content, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, user_id, title, content, created_at, updated_at`,
        userID, req.Title, req.Content, time.Now(),
    ).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating post"})
        return
    }

    c.JSON(http.StatusCreated, post)
}

// @Summary     Get all posts
// @Description Get all posts for the current tenant
// @Tags        posts
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} models.Post "List of posts"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /posts [get]
func GetPosts(c *gin.Context) {
    tenantID := c.GetInt("tenant_id")

    // Get tenant database
    tenantDB, err := database.GetTenantDB(tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // Get all posts
    rows, err := tenantDB.Query(`
        SELECT id, user_id, title, content, created_at, updated_at
        FROM posts
        ORDER BY created_at DESC
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching posts"})
        return
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Content,
            &post.CreatedAt,
            &post.UpdatedAt,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning posts"})
            return
        }
        posts = append(posts, post)
    }

    c.JSON(http.StatusOK, posts)
}

// @Summary     Get a post by ID
// @Description Get a specific post by its ID
// @Tags        posts
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Post ID"
// @Success     200 {object} models.Post "Post details"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     404 {object} map[string]string "Post not found"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /posts/{id} [get]
func GetPost(c *gin.Context) {
    tenantID := c.GetInt("tenant_id")
    postID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Get tenant database
    tenantDB, err := database.GetTenantDB(tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    var post models.Post
    err = tenantDB.QueryRow(`
        SELECT id, user_id, title, content, created_at, updated_at
        FROM posts
        WHERE id = $1`,
        postID,
    ).Scan(
        &post.ID,
        &post.UserID,
        &post.Title,
        &post.Content,
        &post.CreatedAt,
        &post.UpdatedAt,
    )

    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }

    c.JSON(http.StatusOK, post)
} 