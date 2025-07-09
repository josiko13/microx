package service

import (
	"context"
	"fmt"
	"microx/internal/model"
	"microx/internal/repository"
)

type followService struct {
	followRepo   repository.FollowRepository
	userRepo     repository.UserRepository
	timelineRepo repository.TimelineRepository
	tweetRepo    repository.TweetRepository
}

func NewFollowService(
	followRepo repository.FollowRepository,
	userRepo repository.UserRepository,
	timelineRepo repository.TimelineRepository,
	tweetRepo repository.TweetRepository,
) FollowService {
	return &followService{
		followRepo:   followRepo,
		userRepo:     userRepo,
		timelineRepo: timelineRepo,
		tweetRepo:    tweetRepo,
	}
}

// Métodos con receiver (s *followService)
func (s *followService) FollowUser(ctx context.Context, followerID, followingID int64) error {
	// Validaciones básicas
	if followerID == followingID {
		return fmt.Errorf("user cannot follow themselves")
	}

	// Verificar que ambos usuarios existen
	_, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return fmt.Errorf("follower user not found: %w", err)
	}

	_, err = s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return fmt.Errorf("following user not found: %w", err)
	}

	// Verificar si ya existe la relación de follow
	exists, err := s.followRepo.Exists(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("error checking follow relationship: %w", err)
	}

	if exists {
		return fmt.Errorf("user is already following this user")
	}

	// Crear la relación de follow
	follow := &model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	err = s.followRepo.Create(ctx, follow)
	if err != nil {
		return fmt.Errorf("error creating follow relationship: %w", err)
	}

	// Actualizar timeline del follower con tweets del usuario seguido
	if s.timelineRepo != nil {
		err = s.updateFollowerTimeline(ctx, followerID, followingID)
		if err != nil {
			fmt.Printf("Warning: error updating timeline: %v\n", err)
		}
	}

	return nil
}

func (s *followService) UnfollowUser(ctx context.Context, followerID, followingID int64) error {
	// Validaciones básicas
	if followerID == followingID {
		return fmt.Errorf("user cannot unfollow themselves")
	}

	// Verificar que ambos usuarios existen
	_, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return fmt.Errorf("follower user not found: %w", err)
	}

	_, err = s.userRepo.GetByID(ctx, followingID)
	if err != nil {
		return fmt.Errorf("following user not found: %w", err)
	}

	// Verificar si existe la relación de follow
	exists, err := s.followRepo.Exists(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("error checking follow relationship: %w", err)
	}

	if !exists {
		return fmt.Errorf("user is not following this user")
	}

	// Eliminar la relación de follow
	err = s.followRepo.Delete(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("error deleting follow relationship: %w", err)
	}

	// Remover tweets del usuario seguido del timeline del follower
	if s.timelineRepo != nil {
		err = s.removeFromFollowerTimeline(ctx, followerID, followingID)
		if err != nil {
			fmt.Printf("Warning: error removing from timeline: %v\n", err)
		}
	}

	return nil
}

func (s *followService) GetFollowers(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	followers, err := s.followRepo.GetFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting followers: %w", err)
	}

	return followers, nil
}

func (s *followService) GetFollowing(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	following, err := s.followRepo.GetFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting following: %w", err)
	}

	return following, nil
}

func (s *followService) IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error) {
	if followerID == followingID {
		return false, nil
	}

	exists, err := s.followRepo.Exists(ctx, followerID, followingID)
	if err != nil {
		return false, fmt.Errorf("error checking follow relationship: %w", err)
	}

	return exists, nil
}

func (s *followService) updateFollowerTimeline(ctx context.Context, followerID, followingID int64) error {
	tweets, err := s.tweetRepo.GetByUserID(ctx, followingID, 50, 0)
	if err != nil {
		return fmt.Errorf("error getting user tweets: %w", err)
	}

	for _, tweet := range tweets {
		err = s.timelineRepo.AddToTimeline(ctx, followerID, tweet)
		if err != nil {
			return fmt.Errorf("error adding tweet to timeline: %w", err)
		}
	}

	return nil
}

func (s *followService) removeFromFollowerTimeline(ctx context.Context, followerID, followingID int64) error {
	tweets, err := s.tweetRepo.GetByUserID(ctx, followingID, 1000, 0)
	if err != nil {
		return fmt.Errorf("error getting user tweets: %w", err)
	}

	for _, tweet := range tweets {
		err = s.timelineRepo.RemoveFromTimeline(ctx, followerID, tweet.ID)
		if err != nil {
			return fmt.Errorf("error removing tweet from timeline: %w", err)
		}
	}

	return nil
}
