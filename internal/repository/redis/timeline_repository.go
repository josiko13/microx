package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"microx/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type timelineRepository struct {
	client *redis.Client
}

// NewTimelineRepository crea una nueva instancia del repositorio de timeline
func NewTimelineRepository(client *redis.Client) *timelineRepository {
	return &timelineRepository{client: client}
}

// generateTimelineKey genera la clave para el timeline de un usuario
func (r *timelineRepository) generateTimelineKey(userID int64) string {
	return fmt.Sprintf("timeline:%d", userID)
}

// AddToTimeline agrega un tweet al timeline de un usuario
func (r *timelineRepository) AddToTimeline(ctx context.Context, userID int64, tweet *model.TweetWithUser) error {
	key := r.generateTimelineKey(userID)

	// Serializar el tweet a JSON
	tweetJSON, err := json.Marshal(tweet)
	if err != nil {
		return fmt.Errorf("error marshaling tweet: %w", err)
	}

	// Usar el timestamp como score para ordenar por fecha
	score := float64(tweet.CreatedAt.Unix())

	// Agregar al sorted set de Redis
	err = r.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: tweetJSON,
	}).Err()

	if err != nil {
		return fmt.Errorf("error adding tweet to timeline: %w", err)
	}

	// Establecer TTL para el timeline (1 hora por defecto)
	err = r.client.Expire(ctx, key, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("error setting timeline TTL: %w", err)
	}

	return nil
}

// GetTimeline obtiene el timeline de un usuario desde Redis
func (r *timelineRepository) GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
	key := r.generateTimelineKey(userID)

	// Obtener tweets ordenados por score (fecha) descendente
	// ZREVRANGE para obtener los más recientes primero
	result, err := r.client.ZRevRange(ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting timeline from redis: %w", err)
	}

	var tweets []*model.TweetWithUser
	for _, tweetJSON := range result {
		var tweet model.TweetWithUser
		err := json.Unmarshal([]byte(tweetJSON), &tweet)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling tweet: %w", err)
		}
		tweets = append(tweets, &tweet)
	}

	return tweets, nil
}

// RemoveFromTimeline remueve un tweet del timeline
func (r *timelineRepository) RemoveFromTimeline(ctx context.Context, userID int64, tweetID int64) error {
	key := r.generateTimelineKey(userID)

	// Obtener todos los tweets del timeline
	result, err := r.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("error getting timeline for removal: %w", err)
	}

	// Buscar y remover el tweet específico
	for _, tweetJSON := range result {
		var tweet model.TweetWithUser
		err := json.Unmarshal([]byte(tweetJSON), &tweet)
		if err != nil {
			continue // Skip malformed tweets
		}

		if tweet.ID == tweetID {
			err = r.client.ZRem(ctx, key, tweetJSON).Err()
			if err != nil {
				return fmt.Errorf("error removing tweet from timeline: %w", err)
			}
			break
		}
	}

	return nil
}

// InvalidateTimeline invalida el timeline de un usuario (lo elimina del cache)
func (r *timelineRepository) InvalidateTimeline(ctx context.Context, userID int64) error {
	key := r.generateTimelineKey(userID)

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error invalidating timeline: %w", err)
	}

	return nil
}

// AddToMultipleTimelines agrega un tweet a múltiples timelines (para cuando alguien publica)
func (r *timelineRepository) AddToMultipleTimelines(ctx context.Context, followerIDs []int64, tweet *model.TweetWithUser) error {
	// Usar pipeline para operaciones en lote
	pipe := r.client.Pipeline()

	for _, followerID := range followerIDs {
		key := r.generateTimelineKey(followerID)
		score := float64(tweet.CreatedAt.Unix())

		tweetJSON, err := json.Marshal(tweet)
		if err != nil {
			return fmt.Errorf("error marshaling tweet for pipeline: %w", err)
		}

		pipe.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: tweetJSON,
		})
		pipe.Expire(ctx, key, time.Hour)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error executing timeline pipeline: %w", err)
	}

	return nil
}

// GetTimelineSize obtiene el tamaño del timeline de un usuario
func (r *timelineRepository) GetTimelineSize(ctx context.Context, userID int64) (int64, error) {
	key := r.generateTimelineKey(userID)

	size, err := r.client.ZCard(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error getting timeline size: %w", err)
	}

	return size, nil
}
