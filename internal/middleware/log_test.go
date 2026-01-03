package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLogging(t *testing.T) {

	body := strings.NewReader("http://yandex.ru")
	contentType := "text/plain"

	tests := []struct {
		name string // description of this test case
		h    http.HandlerFunc
		want int
	}{
		{
			name: "positive",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) },
			want: http.StatusBadRequest,
		},
		{
			name: "negative",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) },
			want: http.StatusCreated,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := WithLogging(test.h)
			assert.ObjectsAreEqual(test.h, got)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", body)
			r.Header.Add("Content-Type", contentType)
			got.(http.HandlerFunc)(w, r)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, test.want, result.StatusCode, "status code")

		})
	}
}
