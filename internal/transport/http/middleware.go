package http

import (
	"log"
	"net/http"
	"time"
)

type statusWrite struct {
	http.ResponseWriter
	status int
}

func (w *statusWrite) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWrite{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		log.Printf("%s %s => %d (%s)", r.Method, r.URL.Path, sw.status, time.Since(start))
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC: %s %s: %v", r.Method, r.URL.Path, rec)
				writeError(w, http.StatusInternalServerError, "внутреняя ошибка")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
