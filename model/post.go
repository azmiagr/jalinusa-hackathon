package model

import "github.com/google/uuid"

type GetPost struct {
	PostName string `json:"-"`
}

type CreatePostRequest struct {
	PostName   string `json:"post_name" binding:"required"`
	Coordinate string `json:"coordinate" binding:"required"`
	Capacity   int    `json:"capacity" binding:"required"`
}

type CreatePostResponse struct {
	PostID     uuid.UUID `json:"post_id"`
	CreatedBy  string    `json:"created_by"`
	PostCode   string    `json:"post_code"`
	PostName   string    `json:"post_name"`
	Coordinate string    `json:"coordinate"`
	Capacity   int       `json:"capacity"`
	Status     string    `json:"status"`
}

type GetAllPosts struct {
	Posts []PostsResponse
}

type PostsResponse struct {
	PostID     uuid.UUID `json:"post_id"`
	PostCode   string    `json:"post_code"`
	PostName   string    `json:"post_name"`
	KioskCount int       `json:"kiosk_count"`
	Capacity   int       `json:"capacity"`
}
