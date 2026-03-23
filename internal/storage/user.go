package storage

import (
	"context"

	"github.com/effect707/csinder/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{pool: pool}
}

func (s *UserStore) Upsert(ctx context.Context, u *model.User) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO users (faceit_id, nickname, avatar, country, faceit_url, faceit_elo, skill_level, win_rate, kd_ratio, matches, headshot_pct, longest_streak)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (faceit_id) DO UPDATE SET
			nickname = EXCLUDED.nickname,
			avatar = EXCLUDED.avatar,
			country = EXCLUDED.country,
			faceit_url = EXCLUDED.faceit_url,
			faceit_elo = EXCLUDED.faceit_elo,
			skill_level = EXCLUDED.skill_level,
			win_rate = EXCLUDED.win_rate,
			kd_ratio = EXCLUDED.kd_ratio,
			matches = EXCLUDED.matches,
			headshot_pct = EXCLUDED.headshot_pct,
			longest_streak = EXCLUDED.longest_streak,
			updated_at = now()`,
		u.FaceitID, u.Nickname, u.Avatar, u.Country, u.FaceitURL,
		u.FaceitElo, u.SkillLevel, u.WinRate, u.KDRatio, u.Matches,
		u.HeadshotPct, u.LongestStreak,
	)
	return err
}

func (s *UserStore) GetByFaceitID(ctx context.Context, faceitID string) (*model.User, error) {
	return s.scanUser(ctx, "SELECT id, faceit_id, nickname, avatar, country, faceit_url, faceit_elo, skill_level, win_rate, kd_ratio, matches, headshot_pct, longest_streak, created_at, updated_at FROM users WHERE faceit_id = $1", faceitID)
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*model.User, error) {
	return s.scanUser(ctx, "SELECT id, faceit_id, nickname, avatar, country, faceit_url, faceit_elo, skill_level, win_rate, kd_ratio, matches, headshot_pct, longest_streak, created_at, updated_at FROM users WHERE id = $1", id)
}

func (s *UserStore) GetCandidates(ctx context.Context, userID string, elo int, limit int) ([]*model.User, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, faceit_id, nickname, avatar, country, faceit_url, faceit_elo, skill_level, win_rate, kd_ratio, matches, headshot_pct, longest_streak, created_at, updated_at
		FROM users
		WHERE id != $1
		  AND faceit_elo BETWEEN $2 - 300 AND $2 + 300
		  AND id NOT IN (SELECT target_id FROM swipes WHERE swiper_id = $1)
		ORDER BY random()
		LIMIT $3`,
		userID, elo, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		err := rows.Scan(&u.ID, &u.FaceitID, &u.Nickname, &u.Avatar, &u.Country, &u.FaceitURL,
			&u.FaceitElo, &u.SkillLevel, &u.WinRate, &u.KDRatio, &u.Matches,
			&u.HeadshotPct, &u.LongestStreak, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *UserStore) GetStaleUsers(ctx context.Context, limit int) ([]*model.User, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, faceit_id, nickname, avatar, country, faceit_url, faceit_elo, skill_level, win_rate, kd_ratio, matches, headshot_pct, longest_streak, created_at, updated_at
		FROM users
		WHERE updated_at < now() - interval '6 hours'
		ORDER BY updated_at ASC
		LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		err := rows.Scan(&u.ID, &u.FaceitID, &u.Nickname, &u.Avatar, &u.Country, &u.FaceitURL,
			&u.FaceitElo, &u.SkillLevel, &u.WinRate, &u.KDRatio, &u.Matches,
			&u.HeadshotPct, &u.LongestStreak, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *UserStore) scanUser(ctx context.Context, query string, args ...interface{}) (*model.User, error) {
	u := &model.User{}
	err := s.pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.FaceitID, &u.Nickname, &u.Avatar, &u.Country, &u.FaceitURL,
		&u.FaceitElo, &u.SkillLevel, &u.WinRate, &u.KDRatio, &u.Matches,
		&u.HeadshotPct, &u.LongestStreak, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}
