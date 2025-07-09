package service

import (
	"context"
	"fmt"
	"microx/internal/model"
	"microx/internal/repository"
	"strings"
)

type tweetService struct {
	tweetRepo    repository.TweetRepository
	userRepo     repository.UserRepository
	timelineRepo repository.TimelineRepository
	followRepo   repository.FollowRepository
	maxLength    int
}

func NewTweetService(
	tweetRepo repository.TweetRepository,
	userRepo repository.UserRepository,
	timelineRepo repository.TimelineRepository,
	followRepo repository.FollowRepository,
	maxLength int,
) TweetService {
	return &tweetService{
		tweetRepo:    tweetRepo,
		userRepo:     userRepo,
		timelineRepo: timelineRepo,
		followRepo:   followRepo,
		maxLength:    maxLength,
	}
}

// Métodos con receiver (s *tweetService)
func (s *tweetService) CreateTweet(ctx context.Context, userID int64, content string) (*model.TweetResponse, error) {
	// Validar contenido
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("tweet content cannot be empty")
	}

	if len(content) > s.maxLength {
		return nil, fmt.Errorf("tweet content exceeds maximum length of %d characters", s.maxLength)
	}

	// Verificar que el usuario existe
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Crear el tweet
	tweet := &model.Tweet{
		UserID:  userID,
		Content: content,
	}

	err = s.tweetRepo.Create(ctx, tweet)
	if err != nil {
		return nil, fmt.Errorf("error creating tweet: %w", err)
	}

	// Crear tweet con información del usuario para el timeline
	tweetWithUser := &model.TweetWithUser{
		Tweet:    *tweet,
		Username: user.Username,
	}

	// Obtener seguidores del usuario
	followers, err := s.followRepo.GetFollowers(ctx, userID, 1000, 0) // Obtener todos los seguidores
	if err != nil {
		return nil, fmt.Errorf("error getting followers: %w", err)
	}

	// Extraer IDs de seguidores
	var followerIDs []int64
	for _, follower := range followers {
		followerIDs = append(followerIDs, follower.ID)
	}

	// Agregar tweet a los timelines de los seguidores (solo si TimelineRepo está disponible)
	if s.timelineRepo != nil && len(followerIDs) > 0 {
		err = s.timelineRepo.AddToMultipleTimelines(ctx, followerIDs, tweetWithUser)
		if err != nil {
			fmt.Printf("Warning: error adding to timelines: %v\n", err)
		}
	}

	// Crear respuesta
	response := &model.TweetResponse{
		ID:        tweet.ID,
		Content:   tweet.Content,
		UserID:    tweet.UserID,
		Username:  user.Username,
		CreatedAt: tweet.CreatedAt,
	}

	return response, nil
}

func (s *tweetService) GetTweet(ctx context.Context, tweetID int64) (*model.TweetResponse, error) {
	// Obtener tweet con información del usuario (JOIN optimizado)
	tweetWithUser, err := s.tweetRepo.GetByID(ctx, tweetID)
	if err != nil {
		return nil, fmt.Errorf("error getting tweet: %w", err)
	}

	// Crear respuesta directamente desde TweetWithUser
	response := &model.TweetResponse{
		ID:        tweetWithUser.ID,
		Content:   tweetWithUser.Content,
		UserID:    tweetWithUser.UserID,
		Username:  tweetWithUser.Username,
		CreatedAt: tweetWithUser.CreatedAt,
	}

	return response, nil
}

func (s *tweetService) GetUserTweets(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetResponse, error) {
	// Verificar que el usuario existe
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Obtener tweets del usuario con información del usuario (JOIN optimizado)
	tweetsWithUser, err := s.tweetRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting user tweets: %w", err)
	}

	// Convertir a respuesta directamente desde TweetWithUser
	var responses []*model.TweetResponse
	for _, tweetWithUser := range tweetsWithUser {
		response := &model.TweetResponse{
			ID:        tweetWithUser.ID,
			Content:   tweetWithUser.Content,
			UserID:    tweetWithUser.UserID,
			Username:  tweetWithUser.Username,
			CreatedAt: tweetWithUser.CreatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}
