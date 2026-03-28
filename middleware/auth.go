// middleware/auth.go - JWT 인증 미들웨어
// Gin 기반 JWT 인증 미들웨어
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/serateco-dev/comlib-ai/dto"
	"github.com/serateco-dev/comlib-ai/exception"
	"github.com/serateco-dev/comlib-ai/jwt"
)

// AuthMiddleware JWT 인증 미들웨어
func AuthMiddleware(jwtUtil *jwt.JWTUtil, publicPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Public 경로 확인
		if isPublicPath(c.Request.URL.Path, publicPaths) {
			c.Next()
			return
		}

		// Authorization 헤더에서 토큰 추출
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errorResponse := dto.NewErrorResponse(exception.ErrCodeTokenMissing, exception.MsgTokenMissing)
			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		// Bearer 토큰 형식 확인
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errorResponse := dto.NewErrorResponse(exception.ErrCodeTokenInvalid, exception.MsgTokenInvalid)
			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		token := tokenParts[1]

		// 토큰 검증
		if !jwtUtil.IsTokenValid(token) {
			errorResponse := dto.NewErrorResponse(exception.ErrCodeTokenInvalid, exception.MsgTokenInvalid)
			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		// 토큰 만료 확인
		if jwtUtil.IsTokenExpired(token) {
			errorResponse := dto.NewErrorResponse(exception.ErrCodeTokenExpired, exception.MsgTokenExpired)
			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		// 토큰에서 사용자 정보 추출
		claims, err := jwtUtil.ParseToken(token)
		if err != nil {
			errorResponse := dto.NewErrorResponse(exception.ErrCodeTokenInvalid, exception.MsgTokenInvalid)
			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		// 컨텍스트에 사용자 정보 저장
		c.Set("userNo", claims.UserNo)
		c.Set("userId", claims.Subject)
		c.Set("serviceId", claims.ServiceID)
		c.Set("role", claims.Role)
		c.Set("auth", claims.Auth)
		c.Set("nickName", claims.NickName)
		c.Set("provider", claims.Provider)
		c.Set("accessibleApi", claims.AccessibleAPI)
		c.Set("deviceCd", claims.DeviceCD)
		c.Set("deviceStr", claims.DeviceStr)
		c.Set("token", token)

		c.Next()
	}
}

// OptionalAuthMiddleware 선택적 JWT 인증 미들웨어 (토큰이 있으면 검증, 없어도 통과)
func OptionalAuthMiddleware(jwtUtil *jwt.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Bearer 토큰 형식 확인
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]

		// 토큰 검증 (유효하지 않아도 통과)
		if jwtUtil.IsTokenValid(token) && !jwtUtil.IsTokenExpired(token) {
			// 토큰에서 사용자 정보 추출
			if claims, err := jwtUtil.ParseToken(token); err == nil {
				// 컨텍스트에 사용자 정보 저장
				c.Set("userNo", claims.UserNo)
				c.Set("userId", claims.Subject)
				c.Set("serviceId", claims.ServiceID)
				c.Set("role", claims.Role)
				c.Set("auth", claims.Auth)
				c.Set("nickName", claims.NickName)
				c.Set("provider", claims.Provider)
				c.Set("accessibleApi", claims.AccessibleAPI)
				c.Set("deviceCd", claims.DeviceCD)
				c.Set("deviceStr", claims.DeviceStr)
				c.Set("token", token)
			}
		}

		c.Next()
	}
}

// isPublicPath 공개 경로 확인
func isPublicPath(path string, publicPaths []string) bool {
	for _, publicPath := range publicPaths {
		if strings.HasSuffix(publicPath, "/**") {
			// 와일드카드 패턴 처리
			prefix := strings.TrimSuffix(publicPath, "/**")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		} else if strings.HasSuffix(publicPath, "/*") {
			// 단일 레벨 와일드카드 패턴 처리
			prefix := strings.TrimSuffix(publicPath, "/*")
			if strings.HasPrefix(path, prefix) && !strings.Contains(strings.TrimPrefix(path, prefix), "/") {
				return true
			}
		} else if path == publicPath {
			// 정확한 경로 매칭
			return true
		}
	}
	return false
}

// GetUserNo 컨텍스트에서 사용자 번호 추출
func GetUserNo(c *gin.Context) (int64, bool) {
	if userNo, exists := c.Get("userNo"); exists {
		if userNoInt64, ok := userNo.(int64); ok {
			return userNoInt64, true
		}
	}
	return 0, false
}

// GetUserID 컨텍스트에서 사용자 ID 추출
func GetUserID(c *gin.Context) (string, bool) {
	if userID, exists := c.Get("userId"); exists {
		if userIDStr, ok := userID.(string); ok {
			return userIDStr, true
		}
	}
	return "", false
}

// GetRole 컨텍스트에서 역할 추출
func GetRole(c *gin.Context) (string, bool) {
	if role, exists := c.Get("role"); exists {
		if roleStr, ok := role.(string); ok {
			return roleStr, true
		}
	}
	return "", false
}

// GetAuth 컨텍스트에서 권한 목록 추출
func GetAuth(c *gin.Context) ([]string, bool) {
	if auth, exists := c.Get("auth"); exists {
		if authSlice, ok := auth.([]string); ok {
			return authSlice, true
		}
	}
	return nil, false
}

// GetServiceID 컨텍스트에서 서비스 ID 추출
func GetServiceID(c *gin.Context) (string, bool) {
	if serviceID, exists := c.Get("serviceId"); exists {
		if serviceIDStr, ok := serviceID.(string); ok {
			return serviceIDStr, true
		}
	}
	return "", false
}