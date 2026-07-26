package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type sessionContextKey string

const (
	SessionCookieName                     = "session_id"
	sessionIDContextKey sessionContextKey = "session_id"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionCookieName)

		var sessionID string

		if errors.Is(err, http.ErrNoCookie) {
			sessionID = uuid.NewString()
			setSessionCookie(w, SessionCookieName, sessionID)
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			sessionID = cookie.Value
		}

		ctx := context.WithValue(r.Context(), sessionIDContextKey, sessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSessionID(r *http.Request) string {
	sessionID, ok := r.Context().Value(sessionIDContextKey).(string)
	if !ok {
		return ""
	}

	return sessionID
}

func setSessionCookie(w http.ResponseWriter, cookieName, value string) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   int((24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   false, // on prod we should use secure: true :P
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
