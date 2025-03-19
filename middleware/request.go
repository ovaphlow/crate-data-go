package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		url := r.URL.String()

		// Log basic request info
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			// Read body
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = io.ReadAll(r.Body)
				// Restore body
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			log.Printf("%s %s\n", method, url)
			log.Printf("Request Body: %s", string(bodyBytes))
		} else {
			log.Printf("%s %s \n", method, url)
		}

		next.ServeHTTP(w, r)
	})
}
