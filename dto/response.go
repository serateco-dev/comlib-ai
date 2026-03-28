// dto/response.go - ApiResponse.java, ErrorResponse.java 대응
// 표준화된 API 응답 구조체
package dto

import (
	"time"
)

// ApiResponse 성공 응답 구조체
type ApiResponse[T any] struct {
	Data      T         `json:"data,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"requestId,omitempty"`
}

// ErrorResponse 에러 응답 구조체
type ErrorResponse struct {
	RequestID string    `json:"requestId,omitempty"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Path      string    `json:"path,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// NewSuccessResponse 성공 응답 생성
func NewSuccessResponse[T any](data T) *ApiResponse[T] {
	return &ApiResponse[T]{
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewSuccessResponseWithMessage 메시지가 포함된 성공 응답 생성
func NewSuccessResponseWithMessage[T any](data T, message string) *ApiResponse[T] {
	return &ApiResponse[T]{
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewSuccessMessageResponse 데이터 없이 메시지만 포함된 성공 응답 생성
func NewSuccessMessageResponse(message string) *ApiResponse[interface{}] {
	return &ApiResponse[interface{}]{
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse 에러 응답 생성
func NewErrorResponse(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewErrorResponseWithPath 경로가 포함된 에러 응답 생성
func NewErrorResponseWithPath(code, message, path string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Path:      path,
		Timestamp: time.Now(),
	}
}

// SetRequestID 요청 ID 설정
func (r *ApiResponse[T]) SetRequestID(requestID string) *ApiResponse[T] {
	r.RequestID = requestID
	return r
}

// SetRequestID 요청 ID 설정
func (r *ErrorResponse) SetRequestID(requestID string) *ErrorResponse {
	r.RequestID = requestID
	return r
}