package service

import (
	"context"
	"errors"
	"microx/internal/model"
	"testing"
)

func TestTimelineService_GetTimeline(t *testing.T) {
	ctx := context.Background()

	t.Run("éxito desde caché", func(t *testing.T) {
		timelineRepo := &mockTimelineRepo{
			getTimelineFunc: func(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
				return []*model.TweetWithUser{{Tweet: model.Tweet{ID: 1, Content: "cacheado"}}}, nil
			},
		}
		service := NewTimelineService(timelineRepo, &mockTweetRepo{}, &mockUserRepo{}, &mockFollowRepo{})
		resp, err := service.GetTimeline(ctx, 1, 10, 0)
		if err != nil || len(resp) != 1 || resp[0].ID != 1 {
			t.Errorf("esperaba éxito desde caché, obtuve err: %v, resp: %+v", err, resp)
		}
	})

	t.Run("fallback a base de datos", func(t *testing.T) {
		timelineRepo := &mockTimelineRepo{
			getTimelineFunc: func(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
				return nil, errors.New("fallo cache")
			},
		}
		tweetRepo := &mockTweetRepo{
			getTimelineFunc: func(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
				return []*model.TweetWithUser{{Tweet: model.Tweet{ID: 2, Content: "db"}}}, nil
			},
		}
		service := NewTimelineService(timelineRepo, tweetRepo, &mockUserRepo{}, &mockFollowRepo{})
		resp, err := service.GetTimeline(ctx, 1, 10, 0)
		if err != nil || len(resp) != 1 || resp[0].ID != 2 {
			t.Errorf("esperaba fallback a base de datos, obtuve err: %v, resp: %+v", err, resp)
		}
	})

	t.Run("error total", func(t *testing.T) {
		timelineRepo := &mockTimelineRepo{
			getTimelineFunc: func(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
				return nil, errors.New("fallo cache")
			},
		}
		tweetRepo := &mockTweetRepo{
			getTimelineFunc: func(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
				return nil, errors.New("fallo db")
			},
		}
		service := NewTimelineService(timelineRepo, tweetRepo, &mockUserRepo{}, &mockFollowRepo{})
		_, err := service.GetTimeline(ctx, 1, 10, 0)
		if err == nil {
			t.Error("esperaba error total")
		}
	})
}
