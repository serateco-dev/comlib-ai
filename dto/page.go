// dto/page.go - PageRequest.java, PageInfo.java, PagedApiResponse.java 대응
// 페이징 관련 구조체
package dto

import "time"

// PageRequest 페이징 요청 구조체
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"min=1,max=100"`
	Offset   int `json:"offset,omitempty"`
}

// PageInfo 페이징 정보 구조체
type PageInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// PagedApiResponse 페이징된 API 응답 구조체
type PagedApiResponse[T any] struct {
	ApiResponseGeneric[[]T]
	PageInfo PageInfo `json:"pageInfo"`
}

// NewPageRequest 페이징 요청 생성 (기본값 적용)
func NewPageRequest(page, pageSize int) *PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return &PageRequest{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// Init 페이징 요청 초기화
func (p *PageRequest) Init() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	p.Offset = (p.Page - 1) * p.PageSize
}

// GetPage 페이지 번호 반환 (최소값 보장)
func (p *PageRequest) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// GetPageSize 페이지 크기 반환 (기본값 및 최대값 보장)
func (p *PageRequest) GetPageSize() int {
	if p.PageSize < 1 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 오프셋 계산
func (p *PageRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// NewPageInfo 페이징 정보 생성
func NewPageInfo(page, pageSize int, totalCount int64) PageInfo {
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	return PageInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// NewPagedResponse 페이징된 응답 생성
func NewPagedResponse[T any](data []T, pageInfo PageInfo, message ...string) *PagedApiResponse[T] {
	response := &PagedApiResponse[T]{
		ApiResponseGeneric: ApiResponseGeneric[[]T]{
			Success:   true,
			Data:      data,
			Timestamp: time.Now(),
		},
		PageInfo: pageInfo,
	}

	if len(message) > 0 {
		response.ApiResponseGeneric.Message = message[0]
	}

	return response
}