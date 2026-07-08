package model

import "github.com/google/uuid"

type GetUserParam struct {
	UserID   uuid.UUID `json:"-"`
	Username string    `json:"-"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
