package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/effect707/csinder/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStore struct {
	pool *pgxpool.Pool
}

func NewSessionStore(pool *pgxpool.Pool) *SessionStore {
	return &SessionStore{pool: pool}
}

func (s *SessionStore) Create(ctx context.Context, userID string) (*model.Session, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	id := hex.EncodeToString(b)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	_, err := s.pool.Exec(ctx,
		"INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)",
		id, userID, expiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &model.Session{
		ID:        id,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

func (s *SessionStore) Get(ctx context.Context, id string) (*model.Session, error) {
	var sess model.Session
	err := s.pool.QueryRow(ctx,
		"SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = $1 AND expires_at > now()",
		id,
	).Scan(&sess.ID, &sess.UserID, &sess.ExpiresAt, &sess.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *SessionStore) Delete(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM sessions WHERE id = $1", id)
	return err
}

func (s *SessionStore) Cleanup(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM sessions WHERE expires_at < now()")
	return err
}
