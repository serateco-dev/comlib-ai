// middleware/cors.go - CORS 미들웨어
// Cross-Origin Resource Sharing 설정
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 설정
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 기본 CORS 설정
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"X-Gateway-Signature",
			"X-Service-ID",
			"X-Device-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24시간
	}
}

// CORSMiddleware CORS 미들웨어
func CORSMiddleware(config *CORSConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow Origins 설정
		if len(config.AllowOrigins) > 0 {
			if contains(config.AllowOrigins, "*") {
				c.Header("Access-Control-Allow-Origin", "*")
			} else if contains(config.AllowOrigins, origin) {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		// Allow Methods 설정
		if len(config.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowMethods, ", "))
		}

		// Allow Headers 설정
		if len(config.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowHeaders, ", "))
		}

		// Expose Headers 설정
		if len(config.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposeHeaders, ", "))
		}

		// Allow Credentials 설정
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Max Age 설정
		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// OPTIONS 요청 처리
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// contains 문자열 슬라이스에 특정 문자열이 포함되어 있는지 확인
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// joinStrings 문자열 슬라이스를 구분자로 연결
func joinStrings(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}

	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
}