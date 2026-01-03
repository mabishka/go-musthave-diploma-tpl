package compress

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

func TestCompress(t *testing.T) {
	data := "compressed value"
	body := strings.NewReader(data)

	gzipBuffer := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(gzipBuffer)
	gzipWriter.Write([]byte(data))
	gzipWriter.Close()

	deflateBuffer := new(bytes.Buffer)
	deflateWriter := zlib.NewWriter(deflateBuffer)
	deflateWriter.Write([]byte(data))
	deflateWriter.Close()

	tests := []struct {
		name     string // description of this test case
		body     []uint8
		encoding string
		content  string
		wantErr  bool
	}{
		{
			name:     "positive_deflate",
			body:     deflateBuffer.Bytes(),
			encoding: model.CompressTypeDeflate,
			content:  model.ContentTypeJSON,
			wantErr:  false,
		},
		{
			name:     "positive_gzip",
			body:     gzipBuffer.Bytes(),
			encoding: model.CompressTypeGzip,
			content:  model.ContentTypeHTML,
			wantErr:  false,
		},
		{
			name:     "negative_empty_compressType",
			body:     deflateBuffer.Bytes(),
			encoding: "",
			content:  model.ContentTypeJSON,
			wantErr:  true,
		},
		{
			name:     "negative_empty_contentType",
			body:     gzipBuffer.Bytes(),
			encoding: model.CompressTypeDeflate,
			content:  "",
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", body)
			r.Header.Add(model.HeaderAcceptEncoding, test.encoding)
			w.Header().Add(model.HeaderContentType, test.content)
			got := Compress(w, r)

			if test.wantErr {
				assert.Empty(t, got.(*compressResponseWriter).contentEncoding)
				return
			}
			assert.NotEmpty(t, got.(*compressResponseWriter).contentEncoding)
		})
	}
}

func TestDecompress(t *testing.T) {
	data := "compressed value"

	gzipBuffer := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(gzipBuffer)
	gzipWriter.Write([]byte(data))
	gzipWriter.Close()

	deflateBuffer := new(bytes.Buffer)
	deflateWriter := zlib.NewWriter(deflateBuffer)
	deflateWriter.Write([]byte(data))
	deflateWriter.Close()

	tests := []struct {
		name     string // description of this test case
		body     io.Reader
		encoding string
		wantErr  bool
	}{
		{
			name:     "positive_deflate",
			body:     bytes.NewReader(deflateBuffer.Bytes()),
			encoding: model.CompressTypeDeflate,
			wantErr:  false,
		},
		{
			name:     "positive_gzip",
			body:     bytes.NewReader(gzipBuffer.Bytes()),
			encoding: model.CompressTypeGzip,
			wantErr:  false,
		},
		{
			name:     "negative_deflate",
			body:     bytes.NewReader(deflateBuffer.Bytes()),
			encoding: model.CompressTypeGzip,
			wantErr:  true,
		},
		{
			name:     "negative_gzip",
			body:     bytes.NewReader(gzipBuffer.Bytes()),
			encoding: model.CompressTypeDeflate,
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", test.body)
			r.Header.Add(model.HeaderContentEncoding, test.encoding)

			got := Decompress(r)
			assert.NotEmpty(t, got)

			result, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()

			if test.wantErr {
				assert.Empty(t, string(result))
				return
			}
			assert.Equal(t, data, string(result))
		})
	}
}
