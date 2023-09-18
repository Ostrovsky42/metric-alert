package middleware

import (
	"bytes"
	"io"
	"net/http"

	"metric-alert/internal/crypto/hybrid"
	"metric-alert/internal/crypto/symmetric"
	"metric-alert/internal/server/logger"
)

// DecryptorMiddleware представляет промежуточное ПО для декодирования и расшифровки данных HTTP-запросов.
type DecryptorMiddleware struct {
	*hybrid.Decryptor
}

// NewDecryptorMW создает новый экземпляр DecryptorMiddleware с заданным путем к асимметричному ключу.
func NewDecryptorMW(path string) DecryptorMiddleware {
	return DecryptorMiddleware{
		Decryptor: hybrid.NewDecryptor(path),
	}
}

// Decrypt возвращает новый обработчик HTTP, который декодирует и расшифровывает данные запроса перед передачей их следующему обработчику.
func (h DecryptorMiddleware) Decrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.IsNotIncluded() {
			next.ServeHTTP(w, r)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Err(err).Send()
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if h.Symmetric.IsNotCreated() {
			if err := h.SetSymmetric(r.Header.Get(symmetric.HeaderSymmetricKey)); err != nil {
				logger.Log.Err(err).Send()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		data, err = h.DecryptData(data, r.Header.Get(symmetric.HeaderInitializationVector))
		if err != nil {
			logger.Log.Err(err).Send()
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))
		next.ServeHTTP(w, r)
	})
}
