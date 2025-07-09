package service

import (
	"context"
	"fmt"
	"microx/internal/model"
	"microx/internal/repository"
)

type timelineService struct {
	timelineRepo repository.TimelineRepository
	tweetRepo    repository.TweetRepository
	userRepo     repository.UserRepository
	followRepo   repository.FollowRepository
}

// NewTimelineService crea una nueva instancia del servicio de timeline
func NewTimelineService(
	timelineRepo repository.TimelineRepository,
	tweetRepo repository.TweetRepository,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
) TimelineService {
	return &timelineService{
		timelineRepo: timelineRepo,
		tweetRepo:    tweetRepo,
		userRepo:     userRepo,
		followRepo:   followRepo,
	}
}

func (s *timelineService) GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetResponse, error) {
	// Verificar que el usuario existe
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Intentar obtener timeline desde cache (Redis)
	tweets, err := s.timelineRepo.GetTimeline(ctx, userID, limit, offset)
	if err != nil {
		// Si hay error en cache, obtener desde base de datos
		fmt.Printf("Warning: error getting timeline from cache: %v\n", err)
		return s.getTimelineFromDatabase(ctx, userID, limit, offset)
	}

	// Si no hay tweets en cache, obtener desde base de datos
	if len(tweets) == 0 {
		return s.getTimelineFromDatabase(ctx, userID, limit, offset)
	}

	// Convertir a respuesta
	var responses []*model.TweetResponse
	for _, tweet := range tweets {
		response := &model.TweetResponse{
			ID:        tweet.ID,
			Content:   tweet.Content,
			UserID:    tweet.UserID,
			Username:  tweet.Username,
			CreatedAt: tweet.CreatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *timelineService) RefreshTimeline(ctx context.Context, userID int64) error {
	// Verificar que el usuario existe
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Invalidar timeline actual en cache
	err = s.timelineRepo.InvalidateTimeline(ctx, userID)
	if err != nil {
		return fmt.Errorf("error invalidating timeline: %w", err)
	}

	// Obtener timeline desde base de datos
	tweets, err := s.tweetRepo.GetTimeline(ctx, userID, 100, 0) // Obtener primeros 100 tweets
	if err != nil {
		return fmt.Errorf("error getting timeline from database: %w", err)
	}

	// Reconstruir timeline en Redis
	for _, tweet := range tweets {
		err := s.timelineRepo.AddToTimeline(ctx, userID, tweet)
		if err != nil {
			return fmt.Errorf("error adding tweet to timeline: %w", err)
		}
	}

	return nil
}

// getTimelineFromDatabase obtiene el timeline desde la base de datos y lo cachea
func (s *timelineService) getTimelineFromDatabase(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetResponse, error) {
	// Obtener timeline desde base de datos
	tweets, err := s.tweetRepo.GetTimeline(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting timeline from database: %w", err)
	}

	// Cachear tweets en Redis
	for _, tweet := range tweets {
		err := s.timelineRepo.AddToTimeline(ctx, userID, tweet)
		if err != nil {
			// Continuar con el siguiente tweet en lugar de fallar completamente
		}
	}

	// Convertir a respuesta
	var responses []*model.TweetResponse
	for _, tweet := range tweets {
		response := &model.TweetResponse{
			ID:        tweet.ID,
			Content:   tweet.Content,
			UserID:    tweet.UserID,
			Username:  tweet.Username,
			CreatedAt: tweet.CreatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// PreloadAllTimelines carga todos los timelines de todos los usuarios al iniciar la aplicación
func (s *timelineService) PreloadAllTimelines(ctx context.Context) error {
	fmt.Println("🔄 Pre-cargando todos los timelines...")

	// Obtener todos los usuarios
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("error getting all users: %w", err)
	}

	fmt.Printf("📊 Procesando %d usuarios...\n", len(users))

	for _, user := range users {
		// Obtener seguidores del usuario
		followers, err := s.followRepo.GetFollowers(ctx, user.ID, 1000, 0)
		if err != nil {
			fmt.Printf("⚠️ Error getting followers for user %d: %v\n", user.ID, err)
			continue
		}

		if len(followers) == 0 {
			continue // Usuario sin seguidores
		}

		// Obtener tweets del usuario
		tweets, err := s.tweetRepo.GetByUserID(ctx, user.ID, 1000, 0)
		if err != nil {
			fmt.Printf("⚠️ Error getting tweets for user %d: %v\n", user.ID, err)
			continue
		}

		if len(tweets) == 0 {
			continue // Usuario sin tweets
		}

		// Agregar tweets a los timelines de los seguidores
		for _, tweet := range tweets {
			var followerIDs []int64
			for _, follower := range followers {
				followerIDs = append(followerIDs, follower.ID)
			}

			err = s.timelineRepo.AddToMultipleTimelines(ctx, followerIDs, tweet)
			if err != nil {
				fmt.Printf("⚠️ Error adding tweet %d to timelines: %v\n", tweet.ID, err)
			}
		}

		fmt.Printf("✅ Usuario %d: %d tweets agregados a %d timelines\n", user.ID, len(tweets), len(followers))
	}

	fmt.Println("🎉 Pre-carga de timelines completada")
	return nil
}
