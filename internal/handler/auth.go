package handler

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/effect707/csinder/internal/faceit"
	"github.com/effect707/csinder/internal/model"
	"github.com/effect707/csinder/internal/storage"
)

type AuthHandler struct {
	oauth    *faceit.OAuthClient
	api      *faceit.APIClient
	users    *storage.UserStore
	sessions *storage.SessionStore
}

func NewAuthHandler(oauth *faceit.OAuthClient, api *faceit.APIClient, users *storage.UserStore, sessions *storage.SessionStore) *AuthHandler {
	return &AuthHandler{
		oauth:    oauth,
		api:      api,
		users:    users,
		sessions: sessions,
	}
}

func (h *AuthHandler) FaceitLogin(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)
	state := hex.EncodeToString(b)

	verifier, err := faceit.GenerateCodeVerifier()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	challenge := faceit.CodeChallenge(verifier)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "code_verifier",
		Value:    verifier,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, h.oauth.AuthURL(state, challenge), http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	verifierCookie, err := r.Cookie("code_verifier")
	if err != nil || verifierCookie.Value == "" {
		http.Error(w, "missing code verifier", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "oauth_state", Value: "", Path: "/", MaxAge: -1, HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name: "code_verifier", Value: "", Path: "/", MaxAge: -1, HttpOnly: true,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := h.oauth.Exchange(code, verifierCookie.Value)
	if err != nil {
		log.Printf("oauth exchange error: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	claims, err := faceit.DecodeIDToken(token.IDToken)
	if err != nil {
		log.Printf("id_token decode error: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	profile, err := h.api.GetPlayer(claims.GUID)
	if err != nil {
		log.Printf("faceit profile error: %v", err)
		http.Error(w, "failed to fetch profile", http.StatusInternalServerError)
		return
	}

	user := &model.User{
		FaceitID:  claims.GUID,
		Nickname:  profile.Nickname,
		Avatar:    profile.Avatar,
		Country:   profile.Country,
		FaceitURL: profile.FaceitURL,
	}

	if profile.Games.CS2 != nil {
		user.FaceitElo = profile.Games.CS2.FaceitElo
		user.SkillLevel = profile.Games.CS2.SkillLevel
	}

	stats, err := h.api.GetPlayerStats(claims.GUID)
	if err == nil && stats.Lifetime != nil {
		user.WinRate = faceit.ParseFloat(stats.Lifetime["Win Rate %"])
		user.KDRatio = faceit.ParseFloat(stats.Lifetime["Average K/D Ratio"])
		user.Matches = faceit.ParseInt(stats.Lifetime["Matches"])
		user.HeadshotPct = faceit.ParseFloat(stats.Lifetime["Average Headshots %"])
		user.LongestStreak = faceit.ParseInt(stats.Lifetime["Longest Win Streak"])
	}

	if err := h.users.Upsert(r.Context(), user); err != nil {
		log.Printf("user upsert error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	dbUser, err := h.users.GetByFaceitID(r.Context(), claims.GUID)
	if err != nil {
		log.Printf("user fetch error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	sess, err := h.sessions.Create(r.Context(), dbUser.ID)
	if err != nil {
		log.Printf("session create error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		MaxAge:   int(30 * 24 * time.Hour / time.Second),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/swipe", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.sessions.Delete(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
