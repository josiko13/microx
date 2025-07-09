package model

import (
	"time"
)

// Tweet representa un tweet en el sistema
type Tweet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TweetWithUser contiene un tweet con información del usuario
type TweetWithUser struct {
	Tweet
	Username string `json:"username"`
}

// CreateTweetRequest representa la solicitud para crear un tweet
type CreateTweetRequest struct {
	Content string `json:"content" binding:"required,max=280"`
}

// TweetResponse representa la respuesta de un tweet
type TweetResponse struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
