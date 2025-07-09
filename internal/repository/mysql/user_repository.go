package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"microx/internal/model"
	"time"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %d", id)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (username, email, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	// Obtener el ID generado por AUTO_INCREMENT
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	user.ID = id
	return nil
}

func (r *userRepository) GetStats(ctx context.Context, userID int64) (*model.UserStats, error) {
	query := `
		SELECT 
			u.id,
			(SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers_count,
			(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following_count,
			(SELECT COUNT(*) FROM tweets WHERE user_id = u.id) as tweets_count
		FROM users u
		WHERE u.id = ?
	`

	stats := &model.UserStats{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.UserID,
		&stats.FollowersCount,
		&stats.FollowingCount,
		&stats.TweetsCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %d", userID)
		}
		return nil, fmt.Errorf("error getting user stats: %w", err)
	}

	return stats, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at 
		FROM users 
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
