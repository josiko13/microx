package model

import (
	"time"
)

// Follow representa una relación de seguimiento
type Follow struct {
	ID          int64     `json:"id"`
	FollowerID  int64     `json:"follower_id"`
	FollowingID int64     `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// FollowRequest representa la solicitud para seguir a un usuario
type FollowRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// FollowResponse representa la respuesta de una relación de follow
type FollowResponse struct {
	FollowerID  int64     `json:"follower_id"`
	FollowingID int64     `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}
