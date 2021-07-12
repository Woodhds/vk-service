package handlers

import "net/http"

func NotificationHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/event-stream")
		rw.Header().Set("cache-control", "no-cache")
		rw.Header().Set("connection", "keep-alive")
	})
}
