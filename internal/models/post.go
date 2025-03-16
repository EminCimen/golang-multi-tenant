package models

import "time"

// Post represents the post model
type Post struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// CreatePostRequest represents the create post request body
type CreatePostRequest struct {
    Title   string `json:"title" binding:"required" example:"My First Post"`
    Content string `json:"content" binding:"required" example:"This is the content of my first post"`
}

// UpdatePostRequest represents the update post request body
type UpdatePostRequest struct {
    Title   string `json:"title" binding:"required" example:"Updated Post Title"`
    Content string `json:"content" binding:"required" example:"Updated post content"`
} 