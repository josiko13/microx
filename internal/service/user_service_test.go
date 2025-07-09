package service

import (
	"context"
	"errors"
	"microx/internal/model"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Run("creación exitosa", func(t *testing.T) {
		repo := &mockUserRepo{
			createFunc: func(ctx context.Context, user *model.User) error { return nil },
		}
		service := NewUserService(repo)
		user, err := service.CreateUser(context.Background(), "usuario", "mail@mail.com")
		if err != nil || user.Username != "usuario" || user.Email != "mail@mail.com" {
			t.Errorf("esperaba creación exitosa, obtuve err: %v, user: %+v", err, user)
		}
	})

	t.Run("username vacío", func(t *testing.T) {
		service := NewUserService(&mockUserRepo{})
		_, err := service.CreateUser(context.Background(), "", "mail@mail.com")
		if err == nil {
			t.Error("esperaba error por username vacío")
		}
	})

	t.Run("email vacío", func(t *testing.T) {
		service := NewUserService(&mockUserRepo{})
		_, err := service.CreateUser(context.Background(), "usuario", "")
		if err == nil {
			t.Error("esperaba error por email vacío")
		}
	})

	t.Run("error del repositorio", func(t *testing.T) {
		repo := &mockUserRepo{
			createFunc: func(ctx context.Context, user *model.User) error { return errors.New("fallo repo") },
		}
		service := NewUserService(repo)
		_, err := service.CreateUser(context.Background(), "usuario", "mail@mail.com")
		if err == nil || err.Error() != "fallo repo" {
			t.Errorf("esperaba error del repo, obtuve: %v", err)
		}
	})
}
