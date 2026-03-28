// util/time.go - 시간 유틸리티
// 시간 처리 관련 유틸리티 함수
package util

import (
	"time"
)

const (
	// 시간 포맷 상수
	DateFormat         = "2006-01-02"
	TimeFormat         = "15:04:05"
	DateTimeFormat     = "2006-01-02 15:04:05"
	ISO8601Format      = "2006-01-02T15:04:05Z07:00"
	RFC3339Format      = time.RFC3339
	KoreanDateFormat   = "2006년 01월 02일"
	KoreanTimeFormat   = "15시 04분 05초"
	KoreanDateTimeFormat = "2006년 01월 02일 15시 04분 05초"
)

// Now 현재 시간 반환
func Now() time.Time {
	return time.Now()
}

// NowUTC 현재 UTC 시간 반환
func NowUTC() time.Time {
	return time.Now().UTC()
}

// Today 오늘 날짜 (시간은 00:00:00)
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// TodayUTC 오늘 UTC 날짜 (시간은 00:00:00)
func TodayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// StartOfDay 해당 날짜의 시작 시간 (00:00:00)
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 해당 날짜의 끝 시간 (23:59:59.999999999)
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek 해당 주의 시작 (월요일 00:00:00)
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // 일요일
		weekday = 7
	}
	days := weekday - 1
	return StartOfDay(t.AddDate(0, 0, -days))
}

// EndOfWeek 해당 주의 끝 (일요일 23:59:59.999999999)
func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(StartOfWeek(t).AddDate(0, 0, 6))
}

// StartOfMonth 해당 월의 시작 (1일 00:00:00)
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 해당 월의 끝 (마지막 날 23:59:59.999999999)
func EndOfMonth(t time.Time) time.Time {
	return EndOfDay(StartOfMonth(t).AddDate(0, 1, -1))
}

// StartOfYear 해당 년의 시작 (1월 1일 00:00:00)
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 해당 년의 끝 (12월 31일 23:59:59.999999999)
func EndOfYear(t time.Time) time.Time {
	return EndOfDay(time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location()))
}

// FormatDate 날짜 포맷팅
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatTime 시간 포맷팅
func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

// FormatDateTime 날짜시간 포맷팅
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// FormatISO8601 ISO8601 포맷팅
func FormatISO8601(t time.Time) string {
	return t.Format(ISO8601Format)
}

// FormatRFC3339 RFC3339 포맷팅
func FormatRFC3339(t time.Time) string {
	return t.Format(RFC3339Format)
}

// ParseDate 날짜 파싱
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// ParseTime 시간 파싱
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(TimeFormat, timeStr)
}

// ParseDateTime 날짜시간 파싱
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse(DateTimeFormat, dateTimeStr)
}

// ParseISO8601 ISO8601 파싱
func ParseISO8601(iso8601Str string) (time.Time, error) {
	return time.Parse(ISO8601Format, iso8601Str)
}

// ParseRFC3339 RFC3339 파싱
func ParseRFC3339(rfc3339Str string) (time.Time, error) {
	return time.Parse(RFC3339Format, rfc3339Str)
}

// IsToday 오늘인지 확인
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday 어제인지 확인
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}

// IsTomorrow 내일인지 확인
func IsTomorrow(t time.Time) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return t.Year() == tomorrow.Year() && t.Month() == tomorrow.Month() && t.Day() == tomorrow.Day()
}

// IsWeekend 주말인지 확인
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday 평일인지 확인
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// DaysBetween 두 날짜 사이의 일수 계산
func DaysBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}
	return int(end.Sub(start).Hours() / 24)
}

// HoursBetween 두 시간 사이의 시간 계산
func HoursBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}
	return int(end.Sub(start).Hours())
}

// MinutesBetween 두 시간 사이의 분 계산
func MinutesBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}
	return int(end.Sub(start).Minutes())
}

// SecondsBetween 두 시간 사이의 초 계산
func SecondsBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}
	return int(end.Sub(start).Seconds())
}

// AddBusinessDays 영업일 추가 (주말 제외)
func AddBusinessDays(t time.Time, days int) time.Time {
	result := t
	for i := 0; i < days; i++ {
		result = result.AddDate(0, 0, 1)
		for IsWeekend(result) {
			result = result.AddDate(0, 0, 1)
		}
	}
	return result
}

// TimePtr 시간 포인터 생성
func TimePtr(t time.Time) *time.Time {
	return &t
}

// NowPtr 현재 시간 포인터 생성
func NowPtr() *time.Time {
	now := time.Now()
	return &now
}

// IsZero 시간이 제로값인지 확인
func IsZero(t time.Time) bool {
	return t.IsZero()
}

// IsNotZero 시간이 제로값이 아닌지 확인
func IsNotZero(t time.Time) bool {
	return !t.IsZero()
}

// DefaultIfZero 시간이 제로값이면 기본값 반환
func DefaultIfZero(t, defaultTime time.Time) time.Time {
	if t.IsZero() {
		return defaultTime
	}
	return t
}