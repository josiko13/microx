package repository

import (
	"context"
	"microx/internal/model"
)

// UserRepository define las operaciones para usuarios
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	GetStats(ctx context.Context, userID int64) (*model.UserStats, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
}

// TweetRepository define las operaciones para tweets
type TweetRepository interface {
	Create(ctx context.Context, tweet *model.Tweet) error
	GetByID(ctx context.Context, id int64) (*model.TweetWithUser, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error)
	GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error)
}

// FollowRepository define las operaciones para follows
type FollowRepository interface {
	Create(ctx context.Context, follow *model.Follow) error
	Delete(ctx context.Context, followerID, followingID int64) error
	Exists(ctx context.Context, followerID, followingID int64) (bool, error)
	GetFollowers(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error)
	GetFollowing(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error)
}

// TimelineRepository define las operaciones específicas para timeline
type TimelineRepository interface {
	AddToTimeline(ctx context.Context, userID int64, tweet *model.TweetWithUser) error
	GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error)
	RemoveFromTimeline(ctx context.Context, userID int64, tweetID int64) error
	InvalidateTimeline(ctx context.Context, userID int64) error
	AddToMultipleTimelines(ctx context.Context, followerIDs []int64, tweet *model.TweetWithUser) error
}
