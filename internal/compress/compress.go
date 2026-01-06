package compress

import (
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

type ResponseWriter interface {
	http.ResponseWriter
	Close()
}

type compressResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	writer              io.WriteCloser
	contentEncoding     string
}

func (w *compressResponseWriter) Write(b []byte) (int, error) {
	if w.writer != nil {
		return w.writer.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *compressResponseWriter) Close() {
	if w.writer != nil {
		w.writer.Close()
	}
}

func Decompress(r *http.Request) *http.Request {

	if r.Body == nil {
		return r
	}
	decompressType := r.Header.Get(model.HeaderContentEncoding)

	switch decompressType {
	case model.CompressTypeGzip:
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			logger.Log().Error("decompress error",
				zap.Error(err),
				zap.String(model.HeaderContentEncoding, decompressType))
			return r
		}

		r.Body = gz
	case model.CompressTypeDeflate:
		lz, err := zlib.NewReader(r.Body)
		if err != nil {
			logger.Log().Error("decompress error",
				zap.Error(err),
				zap.String(model.HeaderContentEncoding, decompressType))
			return r
		}
		r.Body = lz
	case model.CompressTypeEmpty:
	default:
		logger.Log().Error("decompress error",
			zap.Error(errors.New("unsupport decompress type")),
			zap.String(model.HeaderContentEncoding, decompressType))
	}

	return r
}

func Compress(w http.ResponseWriter, r *http.Request) ResponseWriter {
	cw := &compressResponseWriter{
		ResponseWriter: w,
	}

	content := w.Header().Get(model.HeaderContentType)
	if content != model.ContentTypeJSON && content != model.ContentTypeHTML {
		return cw
	}

	for _, contentEncoding := range r.Header.Values(model.HeaderAcceptEncoding) {
		var compressType string
		compressLevel := 1
		for _, value := range strings.Split(contentEncoding, ",") {
			value = strings.TrimSpace(value)
			if strings.HasPrefix(value, "q=") {
				fmt.Scanf("q=%d", compressLevel)
				continue
			}
			if value != "" {
				compressType = value
			}
		}

		switch compressType {
		case model.CompressTypeGzip:

			gz, err := gzip.NewWriterLevel(w, compressLevel)
			if err != nil {
				logger.Log().Error("compress error",
					zap.Error(err),
					zap.String(model.HeaderContentEncoding, compressType))

				continue
			}
			cw.contentEncoding = compressType
			cw.Header().Set(model.HeaderContentEncoding, compressType)
			cw.writer = gz
			return cw
		case model.CompressTypeDeflate:
			lz, err := zlib.NewWriterLevel(w, compressLevel)
			if err != nil {
				logger.Log().Error("compress error",
					zap.Error(err),
					zap.String(model.HeaderContentEncoding, compressType))

				continue
			}
			cw.contentEncoding = compressType
			cw.Header().Set(model.HeaderContentEncoding, compressType)
			cw.writer = lz
			return cw
		}
	}

	return cw
}
