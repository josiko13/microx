package service

import (
	"context"
	"errors"
	"microx/internal/model"
	"testing"
)

type mockFollowRepo struct {
	getFollowersFunc func(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error)
}

func (m *mockFollowRepo) GetFollowers(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error) {
	if m.getFollowersFunc != nil {
		return m.getFollowersFunc(ctx, userID, limit, offset)
	}
	return []*model.User{}, nil
}
func (m *mockFollowRepo) Create(ctx context.Context, follow *model.Follow) error          { return nil }
func (m *mockFollowRepo) Delete(ctx context.Context, followerID, followingID int64) error { return nil }
func (m *mockFollowRepo) Exists(ctx context.Context, followerID, followingID int64) (bool, error) {
	return false, nil
}
func (m *mockFollowRepo) GetFollowing(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error) {
	return nil, nil
}

func TestTweetService_CreateTweet(t *testing.T) {
	maxLen := 10
	ctx := context.Background()

	t.Run("creación exitosa", func(t *testing.T) {
		tweetRepo := &mockTweetRepo{createFunc: func(ctx context.Context, tweet *model.Tweet) error { return nil }}
		userRepo := &mockUserRepo{getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id, Username: "testuser"}, nil
		}}
		timelineRepo := &mockTimelineRepo{}
		followRepo := &mockFollowRepo{getFollowersFunc: func(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error) {
			return []*model.User{{ID: 2}}, nil
		}}
		service := NewTweetService(tweetRepo, userRepo, timelineRepo, followRepo, maxLen)
		resp, err := service.CreateTweet(ctx, 1, "hola")
		if err != nil || resp.Content != "hola" || resp.UserID != 1 {
			t.Errorf("esperaba creación exitosa, obtuve err: %v, resp: %+v", err, resp)
		}
	})

	t.Run("contenido vacío", func(t *testing.T) {
		service := NewTweetService(&mockTweetRepo{}, &mockUserRepo{}, &mockTimelineRepo{}, &mockFollowRepo{}, maxLen)
		_, err := service.CreateTweet(ctx, 1, "   ")
		if err == nil {
			t.Error("esperaba error por contenido vacío")
		}
	})

	t.Run("contenido demasiado largo", func(t *testing.T) {
		service := NewTweetService(&mockTweetRepo{}, &mockUserRepo{}, &mockTimelineRepo{}, &mockFollowRepo{}, maxLen)
		_, err := service.CreateTweet(ctx, 1, "demasiado largo!")
		if err == nil {
			t.Error("esperaba error por contenido largo")
		}
	})

	t.Run("usuario no existe", func(t *testing.T) {
		userRepo := &mockUserRepo{getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, errors.New("no existe")
		}}
		service := NewTweetService(&mockTweetRepo{}, userRepo, &mockTimelineRepo{}, &mockFollowRepo{}, maxLen)
		_, err := service.CreateTweet(ctx, 1, "hola")
		if err == nil {
			t.Error("esperaba error por usuario no existe")
		}
	})

	t.Run("error al crear tweet", func(t *testing.T) {
		tweetRepo := &mockTweetRepo{createFunc: func(ctx context.Context, tweet *model.Tweet) error { return errors.New("fallo repo") }}
		userRepo := &mockUserRepo{getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id, Username: "testuser"}, nil
		}}
		service := NewTweetService(tweetRepo, userRepo, &mockTimelineRepo{}, &mockFollowRepo{}, maxLen)
		_, err := service.CreateTweet(ctx, 1, "hola")
		if err == nil || err.Error() != "error creating tweet: fallo repo" {
			t.Errorf("esperaba error del repo, obtuve: %v", err)
		}
	})
}
