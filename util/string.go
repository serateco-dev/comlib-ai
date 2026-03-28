// util/string.go - 문자열 유틸리티
// 문자열 처리 관련 유틸리티 함수
package util

import (
	"regexp"
	"strings"
	"unicode"
)

// IsEmpty 문자열이 비어있는지 확인
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty 문자열이 비어있지 않은지 확인
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// HasText 문자열에 텍스트가 있는지 확인 (공백 제외)
func HasText(s string) bool {
	return IsNotEmpty(s)
}

// DefaultIfEmpty 문자열이 비어있으면 기본값 반환
func DefaultIfEmpty(s, defaultValue string) string {
	if IsEmpty(s) {
		return defaultValue
	}
	return s
}

// TrimToEmpty 문자열을 트림하고 nil이면 빈 문자열 반환
func TrimToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}

// TrimToNil 문자열을 트림하고 비어있으면 nil 반환
func TrimToNil(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

// Contains 문자열이 부분 문자열을 포함하는지 확인 (대소문자 구분)
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ContainsIgnoreCase 문자열이 부분 문자열을 포함하는지 확인 (대소문자 무시)
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// StartsWith 문자열이 특정 접두사로 시작하는지 확인
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith 문자열이 특정 접미사로 끝나는지 확인
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// EqualsIgnoreCase 대소문자를 무시하고 문자열 비교
func EqualsIgnoreCase(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

// Capitalize 첫 글자를 대문자로 변환
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uncapitalize 첫 글자를 소문자로 변환
func Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// CamelToSnake 카멜케이스를 스네이크케이스로 변환
func CamelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// SnakeToCamel 스네이크케이스를 카멜케이스로 변환
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	var result strings.Builder
	result.WriteString(parts[0])
	for i := 1; i < len(parts); i++ {
		result.WriteString(Capitalize(parts[i]))
	}
	return result.String()
}

// IsValidEmail 이메일 형식 검증
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// MaskEmail 이메일 마스킹 (개인정보 보호)
func MaskEmail(email string) string {
	if !IsValidEmail(email) {
		return email
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return email
	}

	maskedUsername := string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
	return maskedUsername + "@" + domain
}

// Truncate 문자열을 지정된 길이로 자르기
func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// RemoveWhitespace 모든 공백 문자 제거
func RemoveWhitespace(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, "")
}

// NormalizeWhitespace 연속된 공백을 하나로 정규화
func NormalizeWhitespace(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(s), " ")
}

// Reverse 문자열 뒤집기
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsNumeric 문자열이 숫자인지 확인
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha 문자열이 알파벳인지 확인
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphaNumeric 문자열이 알파벳과 숫자로만 구성되어 있는지 확인
func IsAlphaNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}