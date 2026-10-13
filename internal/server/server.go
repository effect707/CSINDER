package server

import (
	"io/fs"
	"net/http"

	"github.com/effect707/csinder/internal/handler"
	"github.com/effect707/csinder/internal/storage"
)

type Server struct {
	mux *http.ServeMux
}

func New(
	webFS fs.FS,
	sessions *storage.SessionStore,
	users *storage.UserStore,
	authHandler *handler.AuthHandler,
	apiHandler *handler.APIHandler,
	pageHandler *handler.PageHandler,
	wsHandler *handler.WSHandler,
) *Server {
	mux := http.NewServeMux()

	auth := AuthMiddleware(sessions, users)

	staticFS, _ := fs.Sub(webFS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.HandleFunc("GET /{$}", pageHandler.Index)
	mux.HandleFunc("GET /auth/faceit", authHandler.FaceitLogin)
	mux.HandleFunc("GET /auth/callback", authHandler.Callback)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	mux.Handle("GET /swipe", auth(http.HandlerFunc(pageHandler.SwipePage)))
	mux.Handle("GET /profile", auth(http.HandlerFunc(pageHandler.ProfilePage)))
	mux.Handle("GET /matches", auth(http.HandlerFunc(pageHandler.MatchesPage)))

	mux.Handle("GET /api/candidates", auth(http.HandlerFunc(apiHandler.Candidates)))
	mux.Handle("POST /api/swipe", auth(http.HandlerFunc(apiHandler.Swipe)))
	mux.Handle("GET /api/notifications", auth(http.HandlerFunc(apiHandler.Notifications)))
	mux.Handle("POST /api/notifications/seen", auth(http.HandlerFunc(apiHandler.MarkSeen)))

	mux.Handle("GET /ws", auth(http.HandlerFunc(wsHandler.Handle)))

	return &Server{mux: mux}
}

func (s *Server) Handler() http.Handler {
	return LoggingMiddleware(s.mux)
}
