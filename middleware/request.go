package middleware

import (
	"bytes"
	"io"
	"net/http"

	"ovaphlow.com/crate/data/utility"

	"go.uber.org/zap"
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

			utility.ZapLogger.Info("HTTP Request",
				zap.String("method", method),
				zap.String("url", url),
				zap.String("body", string(bodyBytes)),
			)
		} else {
			utility.ZapLogger.Info("HTTP Request",
				zap.String("method", method),
				zap.String("url", url),
			)
		}

		next.ServeHTTP(w, r)
	})
}
