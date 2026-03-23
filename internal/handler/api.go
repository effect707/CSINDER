package handler

import (
	"encoding/json"
	"net/http"

	"github.com/effect707/csinder/internal/auth"
	"github.com/effect707/csinder/internal/service"
	"github.com/effect707/csinder/internal/storage"
)

type APIHandler struct {
	users    *storage.UserStore
	swipes   *storage.SwipeStore
	swipeSvc *service.SwipeService
}

func NewAPIHandler(users *storage.UserStore, swipes *storage.SwipeStore, swipeSvc *service.SwipeService) *APIHandler {
	return &APIHandler{
		users:    users,
		swipes:   swipes,
		swipeSvc: swipeSvc,
	}
}

func (h *APIHandler) Candidates(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	candidates, err := h.users.GetCandidates(r.Context(), user.ID, user.FaceitElo, 10)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidates)
}

func (h *APIHandler) Swipe(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		TargetID string `json:"target_id"`
		Liked    bool   `json:"liked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.TargetID == "" || req.TargetID == user.ID {
		http.Error(w, "invalid target", http.StatusBadRequest)
		return
	}

	if err := h.swipeSvc.Swipe(r.Context(), user.ID, req.TargetID, req.Liked); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *APIHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notifs, err := h.swipes.GetNotifications(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

func (h *APIHandler) MarkSeen(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.swipes.MarkNotificationsSeen(r.Context(), user.ID, req.IDs); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
