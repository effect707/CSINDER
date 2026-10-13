package server

import (
	"log"
	"net/http"
	"time"

	"github.com/effect707/csinder/internal/auth"
	"github.com/effect707/csinder/internal/storage"
)

func AuthMiddleware(sessions *storage.SessionStore, users *storage.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			sess, err := sessions.Get(r.Context(), cookie.Value)
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			user, err := users.GetByID(r.Context(), sess.UserID)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			ctx := auth.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
