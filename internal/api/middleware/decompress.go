package middleware

import (
	"compress/flate"
	"compress/gzip"
	"net/http"
	"strings"
)

func Decompress(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close()
		}
		if strings.Contains(r.Header.Get("Content-Encoding"), "deflate") {
			fl := flate.NewReader(r.Body)
			defer fl.Close()
			r.Body = fl
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
