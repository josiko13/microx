package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"microx/internal/model"
	"time"
)

type tweetRepository struct {
	db *sql.DB
}

// NewTweetRepository crea una nueva instancia del repositorio de tweets
func NewTweetRepository(db *sql.DB) *tweetRepository {
	return &tweetRepository{db: db}
}

func (r *tweetRepository) Create(ctx context.Context, tweet *model.Tweet) error {
	now := time.Now()
	tweet.CreatedAt = now
	tweet.UpdatedAt = now

	query := `
		INSERT INTO tweets (user_id, content, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		tweet.UserID,
		tweet.Content,
		tweet.CreatedAt,
		tweet.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating tweet: %w", err)
	}

	// Obtener el ID generado por AUTO_INCREMENT
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	fmt.Println("ID creado del tweet: ", id)

	tweet.ID = id
	return nil
}

func (r *tweetRepository) GetByID(ctx context.Context, id int64) (*model.TweetWithUser, error) {
	query := `
		SELECT t.id, t.user_id, t.content, t.created_at, t.updated_at, u.username
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`

	tweet := &model.TweetWithUser{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tweet.ID,
		&tweet.UserID,
		&tweet.Content,
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
		&tweet.Username,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tweet not found: %d", id)
		}
		return nil, fmt.Errorf("error getting tweet: %w", err)
	}

	return tweet, nil
}

func (r *tweetRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
	query := `
		SELECT t.id, t.user_id, t.content, t.created_at, t.updated_at, u.username
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = ?
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting user tweets: %w", err)
	}
	defer rows.Close()

	var tweets []*model.TweetWithUser
	for rows.Next() {
		tweet := &model.TweetWithUser{}
		err := rows.Scan(
			&tweet.ID,
			&tweet.UserID,
			&tweet.Content,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
			&tweet.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tweet: %w", err)
		}
		tweets = append(tweets, tweet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tweets: %w", err)
	}

	return tweets, nil
}

func (r *tweetRepository) GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]*model.TweetWithUser, error) {
	query := `
		SELECT t.id, t.user_id, t.content, t.created_at, t.updated_at, u.username
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id IN (
			SELECT following_id 
			FROM follows 
			WHERE follower_id = ?
		)
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting timeline: %w", err)
	}
	defer rows.Close()

	var tweets []*model.TweetWithUser
	for rows.Next() {
		tweet := &model.TweetWithUser{}
		err := rows.Scan(
			&tweet.ID,
			&tweet.UserID,
			&tweet.Content,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
			&tweet.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tweet: %w", err)
		}
		tweets = append(tweets, tweet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating timeline: %w", err)
	}

	return tweets, nil
}
