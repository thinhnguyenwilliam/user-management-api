// user-management-api/internal/middleware/logger.go
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

const (
	HeaderContentType = "Content-Type"

	ContentTypeJSON      = "application/json"
	ContentTypeMultipart = "multipart/form-data"
	ContentTypeForm      = "application/x-www-form-urlencoded"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func LoggerMiddleware() gin.HandlerFunc {
	logger := newHTTPLogger()

	return func(ctx *gin.Context) {
		start := time.Now()

		requestBody := parseRequestBody(ctx)

		writer := wrapResponseWriter(ctx)

		ctx.Next()

		duration := time.Since(start)

		logHTTPEvent(logger, ctx, writer, requestBody, duration)
	}
}

func newHTTPLogger() zerolog.Logger {
	logDir := "internal/logs"

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}

	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logDir, "http-"+today+".log")

	return zerolog.New(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     7,
		Compress:   true,
	}).With().Timestamp().Logger()
}
func wrapResponseWriter(ctx *gin.Context) *CustomResponseWriter {
	writer := &CustomResponseWriter{
		ResponseWriter: ctx.Writer,
		body:           bytes.NewBuffer(nil),
	}
	ctx.Writer = writer
	return writer
}

func parseRequestBody(ctx *gin.Context) map[string]any {
	contentType := ctx.GetHeader(HeaderContentType)

	switch {
	case strings.HasPrefix(contentType, ContentTypeMultipart):
		return parseMultipart(ctx)
	case strings.HasPrefix(contentType, ContentTypeJSON):
		return parseJSONBody(ctx)
	case strings.HasPrefix(contentType, ContentTypeForm):
		return parseFormURLEncoded(ctx)
	default:
		return map[string]any{}
	}
}

func parseJSONBody(ctx *gin.Context) map[string]any {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return map[string]any{}
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var result map[string]any
	_ = json.Unmarshal(bodyBytes, &result)

	return result
}

func parseFormURLEncoded(ctx *gin.Context) map[string]any {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return map[string]any{}
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return map[string]any{}
	}

	result := make(map[string]any)
	for key, vals := range values {
		if len(vals) == 1 {
			result[key] = vals[0]
		} else {
			result[key] = vals
		}
	}

	return result
}

func parseMultipart(ctx *gin.Context) map[string]any {
	result := make(map[string]any)

	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		return result
	}

	form := ctx.Request.MultipartForm
	if form == nil {
		return result
	}

	for key, vals := range form.Value {
		if len(vals) == 1 {
			result[key] = vals[0]
		} else {
			result[key] = vals
		}
	}

	var files []map[string]any
	for field, fileHeaders := range form.File {
		for _, f := range fileHeaders {
			files = append(files, map[string]any{
				"field":        field,
				"filename":     f.Filename,
				"size":         formatFileSize(f.Size),
				"content_type": f.Header.Get(HeaderContentType),
			})
		}
	}

	if len(files) > 0 {
		result["form_files"] = files
	}

	return result
}

func parseResponseBody(writer *CustomResponseWriter) any {
	contentType := writer.Header().Get(HeaderContentType)
	raw := writer.body.String()

	if strings.HasPrefix(contentType, "image/") {
		return "[BINARY DATA]"
	}

	var parsed any
	if json.Unmarshal([]byte(raw), &parsed) == nil {
		return parsed
	}

	return raw
}

func logHTTPEvent(
	logger zerolog.Logger,
	ctx *gin.Context,
	writer *CustomResponseWriter,
	requestBody map[string]any,
	duration time.Duration,
) {
	statusCode := writer.Status()

	logEvent := chooseLogLevel(logger, statusCode)

	logEvent.
		Str("method", ctx.Request.Method).
		Str("path", ctx.Request.URL.Path).
		Str("query", ctx.Request.URL.RawQuery).
		Str("client_ip", ctx.ClientIP()).
		Str("user_agent", ctx.Request.UserAgent()).
		Str("referer", ctx.Request.Referer()).
		Str("protocol", ctx.Request.Proto).
		Str("host", ctx.Request.Host).
		Str("remote_addr", ctx.Request.RemoteAddr).
		Str("request_uri", ctx.Request.RequestURI).
		Int64("content_length", ctx.Request.ContentLength).
		Interface("headers", ctx.Request.Header).
		Interface("request_body", requestBody).
		Int("status_code", statusCode).
		Interface("response_body", parseResponseBody(writer)).
		Int64("duration_ms", duration.Milliseconds()).
		Msg("HTTP Request Log")
}

func chooseLogLevel(logger zerolog.Logger, status int) *zerolog.Event {
	switch {
	case status >= 500:
		return logger.Error()
	case status >= 400:
		return logger.Warn()
	default:
		return logger.Info()
	}
}

func formatFileSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
