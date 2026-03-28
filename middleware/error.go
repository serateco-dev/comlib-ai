// middleware/error.go - 에러 핸들링 미들웨어
// 전역 에러 처리 및 표준화된 에러 응답
package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/serateco-dev/comlib-ai/dto"
	"github.com/serateco-dev/comlib-ai/exception"
)

// ErrorHandlerMiddleware 에러 핸들링 미들웨어
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 에러가 있는 경우 처리
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

// handleError 에러 처리
func handleError(c *gin.Context, err error) {
	var businessErr *exception.BusinessException
	var unauthorizedErr *exception.UnauthorizedException

	switch {
	case errors.As(err, &businessErr):
		// 비즈니스 예외 처리
		errorResponse := dto.NewErrorResponseWithPath(
			businessErr.Code,
			businessErr.Message,
			c.Request.URL.Path,
		)
		c.JSON(businessErr.HTTPStatus, errorResponse)

	case errors.As(err, &unauthorizedErr):
		// 인증 예외 처리
		errorResponse := dto.NewErrorResponseWithPath(
			unauthorizedErr.Code,
			unauthorizedErr.Message,
			c.Request.URL.Path,
		)
		c.JSON(http.StatusUnauthorized, errorResponse)

	default:
		// 기타 예외 처리
		log.Printf("Unhandled error: %v", err)
		errorResponse := dto.NewErrorResponseWithPath(
			exception.ErrCodeInternalError,
			exception.MsgInternalError,
			c.Request.URL.Path,
		)
		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}

// RecoveryMiddleware 패닉 복구 미들웨어
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		
		errorResponse := dto.NewErrorResponseWithPath(
			exception.ErrCodeInternalError,
			exception.MsgInternalError,
			c.Request.URL.Path,
		)
		
		c.JSON(http.StatusInternalServerError, errorResponse)
		c.Abort()
	})
}

// AbortWithError 에러와 함께 요청 중단
func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

// AbortWithBusinessError 비즈니스 에러와 함께 요청 중단
func AbortWithBusinessError(c *gin.Context, code, message string, httpStatus int) {
	err := exception.NewBusinessExceptionWithStatus(code, message, httpStatus)
	AbortWithError(c, err)
}

// AbortWithUnauthorizedError 인증 에러와 함께 요청 중단
func AbortWithUnauthorizedError(c *gin.Context, code, message string) {
	err := exception.NewUnauthorizedException(code, message)
	AbortWithError(c, err)
}