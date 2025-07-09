package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"microx/internal/model"
	"time"
)

type followRepository struct {
	db *sql.DB
}

// NewFollowRepository crea una nueva instancia del repositorio de follows
func NewFollowRepository(db *sql.DB) *followRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(ctx context.Context, follow *model.Follow) error {
	now := time.Now()
	follow.CreatedAt = now

	query := `
		INSERT INTO follows (follower_id, following_id, created_at)
		VALUES (?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		follow.FollowerID,
		follow.FollowingID,
		follow.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating follow: %w", err)
	}

	// Obtener el ID generado por AUTO_INCREMENT
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	follow.ID = id
	return nil
}

func (r *followRepository) Delete(ctx context.Context, followerID, followingID int64) error {
	query := `
		DELETE FROM follows 
		WHERE follower_id = ? AND following_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("error deleting follow: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follow relationship not found")
	}

	return nil
}

func (r *followRepository) Exists(ctx context.Context, followerID, followingID int64) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM follows 
		WHERE follower_id = ? AND following_id = ?
	`

	fmt.Println("query: ", query)
	fmt.Println("followerID: ", followerID)
	fmt.Println("followingID: ", followingID)

	var count int
	err := r.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking follow existence: %w", err)
	}

	return count > 0, nil
}

func (r *followRepository) GetFollowers(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.following_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting followers: %w", err)
	}
	defer rows.Close()

	var followers []*model.User
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
			return nil, fmt.Errorf("error scanning follower: %w", err)
		}
		followers = append(followers, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating followers: %w", err)
	}

	return followers, nil
}

func (r *followRepository) GetFollowing(ctx context.Context, userID int64, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at
		FROM users u
		JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting following: %w", err)
	}
	defer rows.Close()

	var following []*model.User
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
			return nil, fmt.Errorf("error scanning following: %w", err)
		}
		following = append(following, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating following: %w", err)
	}

	return following, nil
}
