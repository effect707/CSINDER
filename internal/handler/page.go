package handler

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/effect707/csinder/internal/auth"
	"github.com/effect707/csinder/internal/storage"
)

type PageHandler struct {
	templates *template.Template
	swipes    *storage.SwipeStore
}

func NewPageHandler(webFS fs.FS, swipes *storage.SwipeStore) *PageHandler {
	tmpl := template.Must(template.ParseFS(webFS, "templates/*.html"))
	return &PageHandler{
		templates: tmpl,
		swipes:    swipes,
	}
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	if cookie != nil && cookie.Value != "" {
		http.Redirect(w, r, "/swipe", http.StatusSeeOther)
		return
	}
	h.templates.ExecuteTemplate(w, "index.html", nil)
}

func (h *PageHandler) SwipePage(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	unseenCount, _ := h.swipes.CountUnseenNotifications(r.Context(), user.ID)
	h.templates.ExecuteTemplate(w, "swipe.html", map[string]interface{}{
		"User":        user,
		"UnseenCount": unseenCount,
	})
}

func (h *PageHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	unseenCount, _ := h.swipes.CountUnseenNotifications(r.Context(), user.ID)
	h.templates.ExecuteTemplate(w, "profile.html", map[string]interface{}{
		"User":        user,
		"UnseenCount": unseenCount,
	})
}

func (h *PageHandler) MatchesPage(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	notifs, _ := h.swipes.GetNotifications(r.Context(), user.ID)
	unseenCount, _ := h.swipes.CountUnseenNotifications(r.Context(), user.ID)
	h.templates.ExecuteTemplate(w, "matches.html", map[string]interface{}{
		"User":          user,
		"Notifications": notifs,
		"UnseenCount":   unseenCount,
	})
}
