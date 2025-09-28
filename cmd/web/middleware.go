package main

import (
	"log/slog"
	"net/http"
	"strings"
)

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
		)
		logger.Info("Request recieved", "method", method, "url", url, "ip", ip)

		h.ServeHTTP(w, r)
	})
}
