package model

import (
	"time"
)

// User representa un usuario en el sistema
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserStats contiene estadísticas del usuario
type UserStats struct {
	UserID         int64 `json:"user_id"`
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
	TweetsCount    int64 `json:"tweets_count"`
}
