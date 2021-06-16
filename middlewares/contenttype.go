package middlewares

import (
	"net/http"
)

func UseContentTypeMiddleWare(handler http.Handler, contentType string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-type", contentType)
		handler.ServeHTTP(rw, r)
	})
}
