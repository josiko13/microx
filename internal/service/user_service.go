package service

import (
	"context"
	"fmt"
	"microx/internal/model"
	"microx/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, email string) (*model.User, error) {
	if username == "" || email == "" {
		return nil, fmt.Errorf("username and email are required")
	}
	user := &model.User{Username: username, Email: email}
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	// Validación básica
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Obtener usuario desde el repositorio
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserStats(ctx context.Context, userID int64) (*model.UserStats, error) {
	// Validación básica
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Verificar que el usuario existe
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Obtener estadísticas del usuario
	stats, err := s.userRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	return stats, nil
}
