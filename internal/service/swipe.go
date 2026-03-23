package service

import (
	"context"
	"time"

	"github.com/effect707/csinder/internal/storage"
	"github.com/effect707/csinder/internal/ws"
)

type SwipeService struct {
	swipes *storage.SwipeStore
	users  *storage.UserStore
	hub    *ws.Hub
}

func NewSwipeService(swipes *storage.SwipeStore, users *storage.UserStore, hub *ws.Hub) *SwipeService {
	return &SwipeService{
		swipes: swipes,
		users:  users,
		hub:    hub,
	}
}

func (s *SwipeService) Swipe(ctx context.Context, swiperID, targetID string, liked bool) error {
	if err := s.swipes.Create(ctx, swiperID, targetID, liked); err != nil {
		return err
	}

	if !liked {
		return nil
	}

	if err := s.swipes.CreateNotification(ctx, targetID, swiperID); err != nil {
		return err
	}

	swiper, err := s.users.GetByID(ctx, swiperID)
	if err != nil {
		return nil
	}

	s.hub.Send(targetID, ws.Message{
		Type: "like_notification",
		Data: map[string]interface{}{
			"from_nickname":   swiper.Nickname,
			"from_avatar":     swiper.Avatar,
			"from_elo":        swiper.FaceitElo,
			"from_faceit_url": swiper.FaceitURL,
			"timestamp":       time.Now().UTC().Format(time.RFC3339),
		},
	})

	return nil
}
