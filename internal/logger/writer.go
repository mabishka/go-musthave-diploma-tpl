package logger

import "net/http"

type ResponseData struct {
	Status  int
	Size    int
	Headers http.Header
}

type LoggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	ResponseData        *ResponseData
}

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size // захватываем размер
	r.ResponseData.Headers = r.Header()
	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode // захватываем код статуса
}
