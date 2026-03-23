package model

import "time"

type User struct {
	ID            string    `json:"id"`
	FaceitID      string    `json:"faceit_id"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	Country       string    `json:"country"`
	FaceitURL     string    `json:"faceit_url"`
	FaceitElo     int       `json:"faceit_elo"`
	SkillLevel    int       `json:"skill_level"`
	WinRate       float64   `json:"win_rate"`
	KDRatio       float64   `json:"kd_ratio"`
	Matches       int       `json:"matches"`
	HeadshotPct   float64   `json:"headshot_pct"`
	LongestStreak int       `json:"longest_streak"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Swipe struct {
	ID        string    `json:"id"`
	SwiperID  string    `json:"swiper_id"`
	TargetID  string    `json:"target_id"`
	Liked     bool      `json:"liked"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	FromUserID string    `json:"from_user_id"`
	Seen       bool      `json:"seen"`
	CreatedAt  time.Time `json:"created_at"`
}

type NotificationWithUser struct {
	Notification
	FromNickname  string `json:"from_nickname"`
	FromAvatar    string `json:"from_avatar"`
	FromElo       int    `json:"from_elo"`
	FromFaceitURL string `json:"from_faceit_url"`
}
