// middleware/logging.go - 로깅 미들웨어
// HTTP 요청/응답 로깅
package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingConfig 로깅 설정
type LoggingConfig struct {
	EnableRequestBody  bool
	EnableResponseBody bool
	SkipPaths          []string
	MaxBodySize        int64
}

// DefaultLoggingConfig 기본 로깅 설정
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		EnableRequestBody:  false,
		EnableResponseBody: false,
		SkipPaths:          []string{"/health", "/metrics"},
		MaxBodySize:        1024, // 1KB
	}
}

// LoggingMiddleware 로깅 미들웨어
func LoggingMiddleware(config *LoggingConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultLoggingConfig()
	}

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Skip paths 확인
		for _, skipPath := range config.SkipPaths {
			if param.Path == skipPath {
				return ""
			}
		}

		return formatLogEntry(param)
	})
}

// DetailedLoggingMiddleware 상세 로깅 미들웨어 (요청/응답 본문 포함)
func DetailedLoggingMiddleware(config *LoggingConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultLoggingConfig()
	}

	return func(c *gin.Context) {
		// Skip paths 확인
		for _, skipPath := range config.SkipPaths {
			if c.Request.URL.Path == skipPath {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		var requestBody []byte
		if config.EnableRequestBody && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, config.MaxBodySize))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Response body 캡처를 위한 writer 래핑
		var responseBody []byte
		if config.EnableResponseBody {
			writer := &responseBodyWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
				maxSize:        config.MaxBodySize,
			}
			c.Writer = writer
			defer func() {
				responseBody = writer.body.Bytes()
			}()
		}

		c.Next()

		// 로그 출력
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logEntry := map[string]interface{}{
			"timestamp":   start.Format(time.RFC3339),
			"latency":     latency.String(),
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"user_agent":  c.Request.UserAgent(),
		}

		if config.EnableRequestBody && len(requestBody) > 0 {
			logEntry["request_body"] = string(requestBody)
		}

		if config.EnableResponseBody && len(responseBody) > 0 {
			logEntry["response_body"] = string(responseBody)
		}

		log.Printf("HTTP Request: %+v", logEntry)
	}
}

// responseBodyWriter 응답 본문을 캡처하는 writer
type responseBodyWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	maxSize int64
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	// 최대 크기 제한
	if w.body.Len() < int(w.maxSize) {
		remaining := int(w.maxSize) - w.body.Len()
		if len(b) > remaining {
			w.body.Write(b[:remaining])
		} else {
			w.body.Write(b)
		}
	}
	return w.ResponseWriter.Write(b)
}

// formatLogEntry 로그 엔트리 포맷팅
func formatLogEntry(param gin.LogFormatterParams) string {
	return log.Sprintf("[%s] %s %s %d %s %s %s\n",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.ClientIP,
		param.Method,
		param.StatusCode,
		param.Latency,
		param.Path,
		param.ErrorMessage,
	)
}