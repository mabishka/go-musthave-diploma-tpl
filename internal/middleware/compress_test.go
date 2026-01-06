package middleware

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestWithCompress(t *testing.T) {
	data := "http://yandex.ru"
	body := strings.NewReader(data)

	gzipBuffer := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(gzipBuffer)
	defer gzipWriter.Close()
	gzipWriter.Write([]byte(data))

	deflateBuffer := new(bytes.Buffer)
	deflateWriter := zlib.NewWriter(deflateBuffer)
	defer deflateWriter.Close()
	deflateWriter.Write([]byte(data))

	tests := []struct {
		name string // description of this test case
		h       http.HandlerFunc
		body    io.Reader
		headers map[string]string
		want    int
	}{
		{
			name: "positive",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) },
			body: body,
			headers: map[string]string{
				model.HeaderContentType: model.ContentTypeJSON,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "positive",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			body: body,
			headers: map[string]string{
				model.HeaderContentType:    model.ContentTypeJSON,
				model.HeaderAcceptEncoding: model.CompressTypeGzip,
			},
			want: http.StatusCreated,
		},
		{
			name: "positive",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			body: body,
			headers: map[string]string{
				model.HeaderContentType:    model.ContentTypeJSON,
				model.HeaderAcceptEncoding: model.CompressTypeGzip,
			},
			want: http.StatusCreated,
		},
		{
			name: "positive",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			body: bytes.NewReader(gzipBuffer.Bytes()),
			headers: map[string]string{
				model.HeaderContentType:     model.ContentTypeJSON,
				model.HeaderAcceptEncoding:  model.CompressTypeGzip,
				model.HeaderContentEncoding: model.CompressTypeGzip,
			},
			want: http.StatusCreated,
		},
		{
			name: "positive",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			body: bytes.NewReader(gzipBuffer.Bytes()),
			headers: map[string]string{
				model.HeaderContentType:     model.ContentTypeJSON,
				model.HeaderAcceptEncoding:  model.CompressTypeDeflate,
				model.HeaderContentEncoding: model.CompressTypeGzip,
			},
			want: http.StatusCreated,
		},
		{
			name: "positive",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			body: bytes.NewReader(deflateBuffer.Bytes()),
			headers: map[string]string{
				model.HeaderContentType:     model.ContentTypeJSON,
				model.HeaderAcceptEncoding:  model.CompressTypeGzip,
				model.HeaderContentEncoding: model.CompressTypeDeflate,
			},
			want: http.StatusCreated,
		},
		{
			name: "positive deflate",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			body: bytes.NewReader(deflateBuffer.Bytes()),
			headers: map[string]string{
				model.HeaderContentType:     model.ContentTypeJSON,
				model.HeaderAcceptEncoding:  model.CompressTypeDeflate,
				model.HeaderContentEncoding: model.CompressTypeDeflate,
			},
			want: http.StatusCreated,
		},
		{
			name: "positive gzip",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			body: bytes.NewReader(deflateBuffer.Bytes()),
			headers: map[string]string{
				model.HeaderContentType:     model.ContentTypeJSON,
				model.HeaderAcceptEncoding:  model.CompressTypeGzip,
				model.HeaderContentEncoding: model.CompressTypeGzip,
			},
			want: http.StatusCreated,
		},
		{
			name: "negative deflate",
			h: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write(deflateBuffer.Bytes())
			},
			body: bytes.NewReader(gzipBuffer.Bytes()),
			headers: map[string]string{
				model.HeaderContentType:     model.ContentTypeJSON,
				model.HeaderAcceptEncoding:  model.CompressTypeDeflate,
				model.HeaderContentEncoding: model.CompressTypeDeflate,
			},
			want: http.StatusCreated,
		},
		{
			name: "negative gzip",
			h: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write(gzipBuffer.Bytes())
			},
			body: bytes.NewReader(deflateBuffer.Bytes()),
			headers: map[string]string{
				model.HeaderContentType:     model.ContentTypeJSON,
				model.HeaderAcceptEncoding:  model.CompressTypeGzip,
				model.HeaderContentEncoding: model.CompressTypeGzip,
			},
			want: http.StatusCreated,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := WithCompress(test.h)
			assert.ObjectsAreEqual(test.h, got)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", test.body)

			for k, v := range test.headers {
				r.Header.Add(k, v)
			}

			got.(http.HandlerFunc)(w, r)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, test.want, result.StatusCode, "status code")
		})
	}
}
