package middleware

import (
	"net/http"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/compress"
)

func WithCompress(h http.Handler) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {

		cr := compress.Decompress(r)
		defer cr.Body.Close()

		cw := compress.Compress(w, r)
		defer cw.Close()

		h.ServeHTTP(cw, cr) // внедряем реализацию http.ResponseWriter
	}

	return http.HandlerFunc(compressFn)
}
