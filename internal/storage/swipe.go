package storage

import (
	"context"

	"github.com/effect707/csinder/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SwipeStore struct {
	pool *pgxpool.Pool
}

func NewSwipeStore(pool *pgxpool.Pool) *SwipeStore {
	return &SwipeStore{pool: pool}
}

func (s *SwipeStore) Create(ctx context.Context, swiperID, targetID string, liked bool) error {
	_, err := s.pool.Exec(ctx,
		"INSERT INTO swipes (swiper_id, target_id, liked) VALUES ($1, $2, $3) ON CONFLICT (swiper_id, target_id) DO NOTHING",
		swiperID, targetID, liked,
	)
	return err
}

func (s *SwipeStore) CreateNotification(ctx context.Context, userID, fromUserID string) error {
	_, err := s.pool.Exec(ctx,
		"INSERT INTO notifications (user_id, from_user_id) VALUES ($1, $2)",
		userID, fromUserID,
	)
	return err
}

func (s *SwipeStore) GetNotifications(ctx context.Context, userID string) ([]*model.NotificationWithUser, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT n.id, n.user_id, n.from_user_id, n.seen, n.created_at,
		       u.nickname, u.avatar, u.faceit_elo, u.faceit_url
		FROM notifications n
		JOIN users u ON u.id = n.from_user_id
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC
		LIMIT 50`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []*model.NotificationWithUser
	for rows.Next() {
		n := &model.NotificationWithUser{}
		err := rows.Scan(&n.ID, &n.UserID, &n.FromUserID, &n.Seen, &n.CreatedAt,
			&n.FromNickname, &n.FromAvatar, &n.FromElo, &n.FromFaceitURL)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, rows.Err()
}

func (s *SwipeStore) MarkNotificationsSeen(ctx context.Context, userID string, ids []string) error {
	_, err := s.pool.Exec(ctx,
		"UPDATE notifications SET seen = true WHERE user_id = $1 AND id = ANY($2)",
		userID, ids,
	)
	return err
}

func (s *SwipeStore) CountUnseenNotifications(ctx context.Context, userID string) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND seen = false",
		userID,
	).Scan(&count)
	return count, err
}
