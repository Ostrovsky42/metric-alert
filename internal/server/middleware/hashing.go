package middleware

import (
	"bytes"
	"io"
	"metric-alert/internal/hasher"
	"net/http"
)

type HashMiddleware struct {
	hb hasher.HashBuilder
}

func NewHashMW(signKey string) HashMiddleware {
	return HashMiddleware{hb: hasher.NewHashGenerator(signKey)}
}

func (h HashMiddleware) Hash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.hb.IsNotActive() {
			next.ServeHTTP(w, r)

			return
		}

		receivedHash := r.Header.Get("HashSHA256")
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		calculatedHash := h.hb.GetHash(data)
		if receivedHash != calculatedHash {
			http.Error(w, "Invalid hash", http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))
		rr := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(rr, r)

		respData := rr.BytesWritten()
		respHash := h.hb.GetHash(respData)
		w.Header().Set("HashSHA256", respHash)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (r *responseRecorder) Write(p []byte) (int, error) {
	r.body.Write(p)
	return r.ResponseWriter.Write(p)
}

func (r *responseRecorder) BytesWritten() []byte {
	return r.body.Bytes()
}
