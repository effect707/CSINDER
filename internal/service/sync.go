package service

import (
	"context"
	"log"
	"time"

	"github.com/effect707/csinder/internal/faceit"
	"github.com/effect707/csinder/internal/storage"
)

type SyncService struct {
	users     *storage.UserStore
	faceitAPI *faceit.APIClient
}

func NewSyncService(users *storage.UserStore, faceitAPI *faceit.APIClient) *SyncService {
	return &SyncService{
		users:     users,
		faceitAPI: faceitAPI,
	}
}

func (s *SyncService) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.syncStaleUsers(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *SyncService) syncStaleUsers(ctx context.Context) {
	users, err := s.users.GetStaleUsers(ctx, 50)
	if err != nil {
		log.Printf("sync: failed to get stale users: %v", err)
		return
	}

	for _, u := range users {
		profile, err := s.faceitAPI.GetPlayer(u.FaceitID)
		if err != nil {
			log.Printf("sync: failed to get profile for %s: %v", u.Nickname, err)
			continue
		}

		stats, err := s.faceitAPI.GetPlayerStats(u.FaceitID)
		if err != nil {
			log.Printf("sync: failed to get stats for %s: %v", u.Nickname, err)
			continue
		}

		u.Nickname = profile.Nickname
		u.Avatar = profile.Avatar
		u.Country = profile.Country
		u.FaceitURL = profile.FaceitURL
		if profile.Games.CS2 != nil {
			u.FaceitElo = profile.Games.CS2.FaceitElo
			u.SkillLevel = profile.Games.CS2.SkillLevel
		}
		if stats.Lifetime != nil {
			u.WinRate = faceit.ParseFloat(stats.Lifetime["Win Rate %"])
			u.KDRatio = faceit.ParseFloat(stats.Lifetime["Average K/D Ratio"])
			u.Matches = faceit.ParseInt(stats.Lifetime["Matches"])
			u.HeadshotPct = faceit.ParseFloat(stats.Lifetime["Average Headshots %"])
			u.LongestStreak = faceit.ParseInt(stats.Lifetime["Longest Win Streak"])
		}

		if err := s.users.Upsert(ctx, u); err != nil {
			log.Printf("sync: failed to update %s: %v", u.Nickname, err)
		}

		time.Sleep(200 * time.Millisecond)
	}
}
