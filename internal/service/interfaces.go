package service

import (
	"context"
	"microx/internal/model"
)

// UserService define las operaciones de negocio para usuarios
type UserService interface {
	CreateUser(ctx context.Context, username, email string) (*model.User, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	GetUserStats(ctx context.Context, userID int64) (*model.UserStats, error)
}

// TweetService define las operaciones de negocio para tweets
type TweetService interface {
	CreateTweet(ctx context.Context, userID int64, content string) (*model.TweetResponse, error)
	GetTweet(ctx context.Context, tweetID int64) (*model.TweetResponse, error)
	GetUserTweets(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetResponse, error)
}

// FollowService define las operaciones de negocio para follows
type FollowService interface {
	FollowUser(ctx context.Context, followerID, followingID int64) error
	UnfollowUser(ctx context.Context, followerID, followingID int64) error
	GetFollowers(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error)
	GetFollowing(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error)
	IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error)
}

// TimelineService define las operaciones de negocio para timeline
type TimelineService interface {
	GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetResponse, error)
	RefreshTimeline(ctx context.Context, userID int64) error
	PreloadAllTimelines(ctx context.Context) error
}
