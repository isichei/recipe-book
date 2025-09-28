package main

import (
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"
	"strings"
)

// Sends WWW-Authenticate header to trigger browser password prompt
func setUnauthed(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func basicAuth(h http.Handler, c creds) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Basic ") {
			setUnauthed(w)
			return
		}

		authString, err := base64.StdEncoding.DecodeString(auth[6:]) // Remove "Basic"
		if err != nil {
			setUnauthed(w)
			return
		}

		credentials := strings.SplitN(string(authString), ":", 2)
		if len(credentials) != 2 {
			setUnauthed(w)
			return
		}

		if !(subtle.ConstantTimeCompare([]byte(credentials[0]), []byte(c.username)) == 1) ||
			!(subtle.ConstantTimeCompare([]byte(credentials[1]), []byte(c.password)) == 1) {
			setUnauthed(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func redirectOldBrowser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_agent := r.Header["User-Agent"][0]
		if strings.Contains(user_agent, "iPad") && strings.Contains(user_agent, "9_3_5") {
			http.Redirect(w, r, "/old", http.StatusSeeOther)
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func logRequest(h http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			header = r.Header.Get("Authorization")
		)
		logger.Info("Request recieved", "method", method, "url", url, "ip", ip, "header", header)

		h.ServeHTTP(w, r)
	})
}
