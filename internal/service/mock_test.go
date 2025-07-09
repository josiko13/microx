package service

import (
	"context"
	"microx/internal/model"
)

type mockUserRepo struct {
	getByIDFunc func(ctx context.Context, id int64) (*model.User, error)
	createFunc  func(ctx context.Context, user *model.User) error
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}
func (m *mockUserRepo) GetStats(ctx context.Context, userID int64) (*model.UserStats, error) {
	return nil, nil
}
func (m *mockUserRepo) GetAllUsers(ctx context.Context) ([]*model.User, error) { return nil, nil }

type mockTimelineRepo struct {
	getTimelineFunc func(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error)
}

func (m *mockTimelineRepo) GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
	if m.getTimelineFunc != nil {
		return m.getTimelineFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (m *mockTimelineRepo) AddToTimeline(ctx context.Context, userID int64, tweet *model.TweetWithUser) error {
	return nil
}
func (m *mockTimelineRepo) RemoveFromTimeline(ctx context.Context, userID int64, tweetID int64) error {
	return nil
}
func (m *mockTimelineRepo) InvalidateTimeline(ctx context.Context, userID int64) error { return nil }
func (m *mockTimelineRepo) AddToMultipleTimelines(ctx context.Context, followerIDs []int64, tweet *model.TweetWithUser) error {
	return nil
}

type mockTweetRepo struct {
	getTimelineFunc func(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error)
	createFunc      func(ctx context.Context, tweet *model.Tweet) error
}

func (m *mockTweetRepo) GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
	if m.getTimelineFunc != nil {
		return m.getTimelineFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}
func (m *mockTweetRepo) Create(ctx context.Context, tweet *model.Tweet) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, tweet)
	}
	return nil
}
func (m *mockTweetRepo) GetByID(ctx context.Context, id int64) (*model.TweetWithUser, error) {
	return nil, nil
}
func (m *mockTweetRepo) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
	return nil, nil
}
