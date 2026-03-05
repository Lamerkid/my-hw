package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("%s [%s] %s %s %s %d %d \"%s\"\n",
			strings.Split(r.RemoteAddr, ":")[0],
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			http.StatusOK,
			time.Since(start),
			strings.Split(r.UserAgent(), " ")[0])
	})
}
