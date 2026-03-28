// exception/exceptions.go - BusinessException.java, UnauthorizedException.java 대응
// 비즈니스 예외 및 인증 예외 정의
package exception

import (
	"fmt"
	"net/http"
)

// BusinessException 비즈니스 예외
type BusinessException struct {
	Code       string
	Message    string
	HTTPStatus int
	Cause      error
}

// Error Error 인터페이스 구현
func (e *BusinessException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 원본 에러 반환
func (e *BusinessException) Unwrap() error {
	return e.Cause
}

// NewBusinessException 비즈니스 예외 생성 (기본 HTTP 500)
func NewBusinessException(code, message string) *BusinessException {
	return &BusinessException{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// NewBusinessExceptionWithStatus HTTP 상태 코드를 지정하는 비즈니스 예외 생성
func NewBusinessExceptionWithStatus(code, message string, httpStatus int) *BusinessException {
	return &BusinessException{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// NewBusinessExceptionWithCause 원인을 포함한 비즈니스 예외 생성 (기본 HTTP 500)
func NewBusinessExceptionWithCause(code, message string, cause error) *BusinessException {
	return &BusinessException{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Cause:      cause,
	}
}

// NewBusinessExceptionWithStatusAndCause HTTP 상태 코드와 원인을 포함한 비즈니스 예외 생성
func NewBusinessExceptionWithStatusAndCause(code, message string, httpStatus int, cause error) *BusinessException {
	return &BusinessException{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Cause:      cause,
	}
}

// UnauthorizedException 401 Unauthorized 전용 예외
type UnauthorizedException struct {
	Code    string
	Message string
	Cause   error
}

// Error Error 인터페이스 구현
func (e *UnauthorizedException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 원본 에러 반환
func (e *UnauthorizedException) Unwrap() error {
	return e.Cause
}

// NewUnauthorizedException 인증 예외 생성
func NewUnauthorizedException(code, message string) *UnauthorizedException {
	return &UnauthorizedException{
		Code:    code,
		Message: message,
	}
}

// NewUnauthorizedExceptionWithCause 원인을 포함한 인증 예외 생성
func NewUnauthorizedExceptionWithCause(code, message string, cause error) *UnauthorizedException {
	return &UnauthorizedException{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// 공통 에러 코드 상수
const (
	// 인증 관련
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeTokenExpired      = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid      = "TOKEN_INVALID"
	ErrCodeTokenMissing      = "TOKEN_MISSING"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"

	// 비즈니스 로직 관련
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND"
	ErrCodeDuplicateResource = "DUPLICATE_RESOURCE"
	ErrCodeInternalError    = "INTERNAL_ERROR"

	// 데이터베이스 관련
	ErrCodeDatabaseError     = "DATABASE_ERROR"
	ErrCodeConnectionFailed  = "CONNECTION_FAILED"
	ErrCodeTransactionFailed = "TRANSACTION_FAILED"
)

// 공통 에러 메시지
const (
	MsgUnauthorized      = "인증이 필요합니다"
	MsgTokenExpired      = "토큰이 만료되었습니다"
	MsgTokenInvalid      = "유효하지 않은 토큰입니다"
	MsgTokenMissing      = "토큰이 없습니다"
	MsgInvalidCredentials = "잘못된 인증 정보입니다"
	MsgValidationFailed  = "입력값 검증에 실패했습니다"
	MsgResourceNotFound  = "요청한 리소스를 찾을 수 없습니다"
	MsgDuplicateResource = "이미 존재하는 리소스입니다"
	MsgInternalError     = "내부 서버 오류가 발생했습니다"
	MsgDatabaseError     = "데이터베이스 오류가 발생했습니다"
	MsgConnectionFailed  = "연결에 실패했습니다"
)