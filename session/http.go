package session

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
)

type ctxKey int

const (
	// CookieName is the name to store the cookie in the browser with.
	CookieName                = "_gitpods_session"
	CookieUserID       ctxKey = iota
	CookieUserUsername ctxKey = iota
)

var (
	errUnauthorized = []*jsonapi.ErrorObject{{
		Title:  http.StatusText(http.StatusUnauthorized),
		Detail: "Your Cookie is not valid",
		Status: fmt.Sprintf("%d", http.StatusUnauthorized),
	}}
)

// Authorized users will have a user information in the next handlers.
func Authorized(s Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(CookieName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			if cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			if time.Now().Before(cookie.Expires) {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			session, err := s.FindSession(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				jsonapi.MarshalErrors(w, errUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CookieUserID, session.User.ID)
			ctx = context.WithValue(ctx, CookieUserUsername, session.User.Username)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetSessionUser from the http.Request
func GetSessionUser(ctx context.Context) *User {
	return &User{
		ID:       ctx.Value(CookieUserID).(string),
		Username: ctx.Value(CookieUserUsername).(string),
	}
}
