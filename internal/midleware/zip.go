package midleware

import (
	"compress/gzip"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"strings"
)

type ZipMiddleware struct {
	zContentTypes []string
	log           zerolog.Logger
	gzipW         *gzip.Writer
}

func NewZipMiddleware(log zerolog.Logger, level int) ZipMiddleware {
	gzW, err := gzip.NewWriterLevel(nil, level)
	if err != nil {
		log.Fatal().Err(err).Msg("err create gzip writer")
	}

	return ZipMiddleware{
		log:   log,
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

func (z *ZipMiddleware) Zip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if z.isNeedZipped(r) {
			defer z.gzipW.Close()

			w.Header().Set("Content-Encoding", "gzip")
			z.gzipW.Reset(w)

			next.ServeHTTP(&gzipWriter{ResponseWriter: w, gzipW: z.gzipW}, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (z *ZipMiddleware) UnZip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)

			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			z.log.Err(err).Msg("err UnZip request body")

			return
		}
		defer gz.Close()

		r.Body = io.NopCloser(gz)

		next.ServeHTTP(w, r)
	})
}

func (z *ZipMiddleware) isNeedZipped(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}
