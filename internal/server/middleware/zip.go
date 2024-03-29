package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"metric-alert/internal/server/logger"
)

type ZipMiddleware struct {
	gzipW *gzip.Writer
}

func NewZipMiddleware(level int) ZipMiddleware {
	gzW, err := gzip.NewWriterLevel(nil, level)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err create gzip writer")
	}

	return ZipMiddleware{
		gzipW: gzW,
	}
}

type gzipWriter struct {
	http.ResponseWriter
	gzipW *gzip.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.gzipW.Write(b)
}

// Zip добавляет сжатие gzip к обработке HTTP-запросов.
func (z *ZipMiddleware) Zip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if z.isNeedZipped(r) {
			defer z.gzipW.Flush()
			defer z.gzipW.Close()

			w.Header().Set("Content-Encoding", "gzip")
			z.gzipW.Reset(w)

			next.ServeHTTP(&gzipWriter{ResponseWriter: w, gzipW: z.gzipW}, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// UnZip добавляет распаковку данных из gzip к обработке HTTP-запросов.
func (z *ZipMiddleware) UnZip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)

			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Log.Err(err).Msg("err UnZip request body")

			return
		}
		defer gz.Close()

		r.Body = io.NopCloser(gz)

		next.ServeHTTP(w, r)
	})
}

// isNeedZipped проверяет, нужно ли сжимать данные на основе заголовков запроса.
func (z *ZipMiddleware) isNeedZipped(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}
