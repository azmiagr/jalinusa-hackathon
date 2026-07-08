package model

import (
	"time"

	"github.com/google/uuid"
)

type GetPost struct {
	PostID   uuid.UUID `json:"-"`
	PostName string    `json:"-"`
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

type GetPostResponse struct {
	PostID   uuid.UUID `json:"post_id"`
	PostCode string    `json:"post_code"`
	PostName string    `json:"post_name"`
}

type BindPostRequest struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	PostID     uuid.UUID `json:"post_id"`
	DeviceName string    `json:"device_name"`
}

type BindPostResponse struct {
	PostID     uuid.UUID `json:"post_id"`
	Status     string    `json:"status"`
	DeviceName string    `json:"device_name"`
	BoundBy    string    `json:"bound_by"`
	BountAt    time.Time `json:"bount_at"`
}
