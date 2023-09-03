// Пакет middleware предоставляет промежуточное программное обеспечение для обработки HTTP-запросов.
package middleware

import (
	"bytes"
	"io"
	"net/http"

	"metric-alert/internal/hasher"
)

// HashMiddleware представляет структуру middleware, выполняющего хеширование.
type HashMiddleware struct {
	hb hasher.HashBuilder
}

// NewHashMW создает новый экземпляр HashMiddleware.
func NewHashMW(signKey string) HashMiddleware {
	return HashMiddleware{hb: hasher.NewHashGenerator(signKey)}
}

// Hash если передан хэшер с ключом сравнивает хэш из заголовка с с хэшом  полученного тела HTTP-запроса с использованием переданного ключа.
// Так же возвращает в заголовке хэш ответа.
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
